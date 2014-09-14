// Receives docker webhooks and triggers updates for them
package main

import (
	"encoding/json"
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
			return
		}
		log.Println("Received webhook for", hook.Repository.RepoName)
		res.Write([]byte("OK"))
		go rollingUpdate([]string{hook.Repository.RepoName})
	})
	http.ListenAndServe(":6600", nil)
}
