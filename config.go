package main

import (
	"encoding/json"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"github.com/ogier/pflag"
	"io/ioutil"
	"log"
	"os"
)

var configFilePath = pflag.String(
	"config",
	"./containers.json",
	"The container list to maintain",
)

type Config []*docker.Container

var clusterConfig Config

// Read the config file into clusterConfig
func parseConfig() bool {
	data, err := ioutil.ReadFile(*configFilePath)
	if err != nil {
		log.Println("Couldn't read the config file!", err)
		return false
	}
	err = json.Unmarshal(data, &clusterConfig)
	if err != nil {
		log.Println("Couldn't parse config!", err)
	}
	return true
}

// Attempts to start a container using its ID
func bringUpContainer(cfgCon *docker.Container) {
	if len(cfgCon.ID) < 10 {
		log.Println("Container was not found, creating it...")
		createContainer(cfgCon)
	}
	cli := newClient()
	err := cli.StartContainer(cfgCon.ID, cfgCon.HostConfig)
	if err != nil {
		log.Println("Error starting the container!")
		log.Println(err)
	}
}

// Get the list of running containers
func listContainers() []*docker.Container {
	cli := newClient()
	opts := docker.ListContainersOptions{
		All: true,
	}
	list, err := cli.ListContainers(opts)
	if err != nil {
		panic(err)
	}
	cons := []*docker.Container{}
	for _, con := range list {
		data, err := cli.InspectContainer(con.ID)
		if err != nil {
			panic(err)
		}
		cons = append(cons, data)
	}
	return cons
}

// Create a new config file based on the current state if none is given
func createConfig() {
	list := listContainers()
	if len(list) < 1 {
		log.Println("No containers found! Exiting...")
		os.Exit(1)
	}
	js, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		panic(err)
	}
	log.Println(string(js))
	fi, err := os.Create(*configFilePath)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(fi, string(js))
}

// Reads new container information into config after an update
func updateConfig() {
	log.Println("Updating the config...")
	list := listContainers()
	up := []*docker.Container{}
	for _, cfgCon := range clusterConfig {
		for _, con := range list {
			if cfgCon.Name == con.Name {
				log.Println("Adding", con.Name, "to config")
				up = append(up, con)
			}
		}
	}
	js, err := json.MarshalIndent(up, "", "  ")
	if err != nil {
		panic(err)
	}
	log.Println("Updated the config", string(js))
	fi, err := os.Create(*configFilePath)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(fi, string(js))
}
