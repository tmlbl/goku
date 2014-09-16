// shgod maintains a state of docker containers defined in a config file.
// It must be run as a user that has access to the docker endpoint used
package main

import (
	"github.com/fsouza/go-dockerclient"
	"github.com/ogier/pflag"
	//"gopkg.in/yaml.v1"
	"log"
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
	pflag.Parse()
	initLogger()
	for _, arg := range os.Args {
		if arg == "init" {
			createConfig()
		}
		if arg == "update" {
			rollingUpdate(os.Args[2:])
			os.Exit(0)
		}
	}
	success := parseConfig()
	if !success {
		log.Println("Couldn't read the config files." +
			" Run shgod init to create them.")
		os.Exit(1)
	}
	log.Println("Using config with containers:")
	for _, con := range clusterConfig {
		log.Println(con.Name)
	}
	go heartbeat()
	serve()
}

func heartbeat() {
	for {
		parseConfig()
		time.Sleep(*interval)
		log.Println(time.Now().String())
		checkState()
	}
}

func newClient() *docker.Client {
	cli, err := docker.NewClient(*endpoint)
	if err != nil {
		log.Fatalln("Failed to connect to docker!", err)
	}
	return cli
}
