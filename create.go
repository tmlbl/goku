// If a container in the config file has been removed, create a new one
package main

import (
	"github.com/fsouza/go-dockerclient"
	"log"
)

func createContainer(con *docker.Container) {
	// Delete a conflicting container if it exists
	removeContainer(con)
	opts := docker.CreateContainerOptions{
		Name:   con.Name,
		Config: con.Config,
	}
	cli := newClient()
	newCon, err := cli.CreateContainer(opts)
	if err != nil {
		logErr(err)
	}
	log.Println("Successfully created container", con)
	// Start the container
	cli.StartContainer(newCon.ID, con.HostConfig)
}

func removeContainer(con *docker.Container) {
	cli := newClient()
	opts := docker.RemoveContainerOptions{
		ID:    con.ID,
		Force: true,
	}
	cli.RemoveContainer(opts)
}
