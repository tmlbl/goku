// shgod maintains a state of docker containers defined in a config file.
// It must be run as a user that has access to the docker endpoint used
package main

import (
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"github.com/ogier/pflag"
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
	pflag.Parse()
	for _, arg := range os.Args {
		if arg == "init" {
			createConfig()
		}
	}
	success := parseConfig()
	if !success {
		fmt.Println("Couldn't read the config files. Run shgod init to create them")
		os.Exit(1)
	}
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("OK\n"))
	})
	go heartbeat()
	http.ListenAndServe(":6600", nil)
}

func heartbeat() {
	for {
		time.Sleep(*interval)
		fmt.Println("Heartbeat")
		list := listContainers()
		for _, cfgCon := range clusterConfig {
			found := false
			for _, con := range list {
				if con.ID == cfgCon.ID {
					checkState(con, cfgCon)
					found = true
				}
			}
			if !found {
				checkState(docker.APIContainers{}, cfgCon)
			}
		}
	}
}

func newClient() *docker.Client {
	cli, err := docker.NewClient(*endpoint)
	if err != nil {
		fmt.Println("Failed to connect to docker! Error:")
		fmt.Println(err)
	}
	return cli
}
