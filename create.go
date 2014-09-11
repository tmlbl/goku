// If a container in the config file has been removed, create a new one
package main

import (
	"fmt"
	"github.com/fsouza/go-dockerclient"
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
	fmt.Println("Successfully created container", con)
	// Update the clusterConfig
	// TODO: Instead of reinitializing clusterConfig, write only the changes
	createConfig()
	parseConfig()
	// Start the container
	cli.StartContainer(newCon.ID, &docker.HostConfig{})
}
