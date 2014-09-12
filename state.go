package main

import (
	"github.com/fsouza/go-dockerclient"
	"log"
)

// Check the current state against the config state
func checkState(con *docker.Container, cfgCon *docker.Container) {
	up := con.State.Running
	if !up {
		log.Println(cfgCon.Name, "is down, bringing it up...")
		bringUpContainer(cfgCon)
	} else {
		log.Println(cfgCon.Name, "is up")
	}
}
