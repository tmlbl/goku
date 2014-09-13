// shgod maintains a state of docker containers defined in a config file.
// It must be run as a user that has access to the docker endpoint used
package main

import (
	"github.com/fsouza/go-dockerclient"
	"github.com/ogier/pflag"
	"gopkg.in/yaml.v1"
	"log"
	"net/http"
	"os"
	"time"
)

var endpoint = pflag.String(
	"endpoint",
	"unix://var/run/docker.sock",
	"The docker endpoint",
)

var interval = pflag.Duration(
	"heartbeat",
	900*time.Millisecond,
	"The heartbeat interval for container health checks",
)

func main() {
	initLogger()
	pflag.Parse()
	for _, arg := range os.Args {
		if arg == "init" {
			createConfig()
		}
		if arg == "update" {
			rollingUpdate(os.Args[2:])
		}
	}
	success := parseConfig()
	if !success {
		log.Println("Couldn't read the config files." +
			" Run shgod init to create them.")
		os.Exit(1)
	}
	y, err := yaml.Marshal(clusterConfig)
	if err != nil {
		log.Panicln("Malformed yaml in config!", err)
	}
	s := string(y)
	log.Println(s)
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		res.Write(y)
	})
	go heartbeat()
	http.ListenAndServe(":6600", nil)
}

func heartbeat() {
	for {
		time.Sleep(*interval)
		log.Println(time.Now().String())
		list := listContainers()
		for _, cfgCon := range clusterConfig {
			found := false
			for _, con := range list {
				if con.Name == cfgCon.Name {
					checkState(con, cfgCon)
					found = true
				}
			}
			if !found {
				log.Println("No container found for", cfgCon.Name)
				createContainer(cfgCon)
				updateConfig()
			}
		}
	}
}

func newClient() *docker.Client {
	cli, err := docker.NewClient(*endpoint)
	if err != nil {
		log.Fatalln("Failed to connect to docker!", err)
	}
	return cli
}
