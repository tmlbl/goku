// If a container in the config file has been removed, create a new one
package main

import (
	"github.com/fsouza/go-dockerclient"
	"log"
)

func createContainer(con *docker.Container) {
	opts := docker.CreateContainerOptions{
		Name:   con.Name,
		Config: con.Config,
	}
	cli := newClient()
	newCon, err := cli.CreateContainer(opts)
	if err != nil {
		panic(err)
	}
	log.Println("Successfully created container", con)
	// Start the container
	cli.StartContainer(newCon.ID, con.HostConfig)
}
