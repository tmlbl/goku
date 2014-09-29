// Receives docker webhooks and triggers updates for them
package main

import (
	"encoding/json"
	"fmt"
	"github.com/ogier/pflag"
	"io/ioutil"
	"log"
	"net/http"
)

type Webhook struct {
	Repository Repository `json:"repository"`
}

type Repository struct {
	RepoName string `json:"repo_name"`
}

var port = pflag.Int32(
	"port",
	6600,
	"Port to listen for webhooks on",
)

func serve() {
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			panic(err)
		}
		hook := &Webhook{}
		hook.Repository = Repository{}
		err = json.Unmarshal(body, hook)
		if err != nil {
			logErr(err)
			res.Write([]byte("OK"))
			return
		}
		log.Println("Received webhook for", hook.Repository.RepoName)
		res.Write([]byte("OK"))
		go rollingUpdate([]string{hook.Repository.RepoName})
	})
	portcfg := fmt.Sprintf(":%d", *port)
	log.Println("Listening at", portcfg)
	http.ListenAndServe(portcfg, nil)
}
