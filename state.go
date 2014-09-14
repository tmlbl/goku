package main

import (
	"github.com/fsouza/go-dockerclient"
	"log"
)

// Check the current state against the config state
// for a single container
func checkContainerState(con *docker.Container, cfgCon *docker.Container) {
	up := con.State.Running
	if !up {
		log.Println(cfgCon.Name, "is down, bringing it up...")
		bringUpContainer(cfgCon)
	}
}

// List the containers and compare them to the config container list
func checkState() {
	list := listContainers()
	allup := true
	for _, cfgCon := range clusterConfig {
		found := false
		for _, con := range list {
			if con.Name == cfgCon.Name {
				if !con.State.Running {
					allup = false
				}
				checkContainerState(con, cfgCon)
				found = true
			}
		}
		if !found {
			log.Println("No container found for", cfgCon.Name)
			createContainer(cfgCon)
			updateConfig()
		}
	}
	if allup {
		log.Println("All containers are up")
	}
}
