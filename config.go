package main

import (
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"github.com/ogier/pflag"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"log"
	"os"
)

var configFilePath = pflag.String(
	"config",
	"./containers.yml",
	"The container list to maintain",
)

type Config []*docker.Container

var clusterConfig Config

// Read the config file into clusterConfig
func parseConfig() bool {
	data, err := ioutil.ReadFile(*configFilePath)
	if err != nil {
		logErr(err)
		return false
	}
	clusterConfig = Config{}
	err = yaml.Unmarshal(data, &clusterConfig)
	if err != nil {
		logErr(err)
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
		logErr(err)
	}
}

// Get a list of inspected containers
func listContainers() []*docker.Container {
	cli := newClient()
	opts := docker.ListContainersOptions{
		All: true,
	}
	list, err := cli.ListContainers(opts)
	if err != nil {
		logErr(err)
	}
	cons := []*docker.Container{}
	for _, con := range list {
		data, err := cli.InspectContainer(con.ID)
		if err == nil {
			cons = append(cons, data)
		}
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
	y, err := yaml.Marshal(list)
	if err != nil {
		logErr(err)
	}
	s := string(y)
	log.Println(s)
	fi, err := os.Create(*configFilePath)
	if err != nil {
		logErr(err)
	}
	fmt.Fprintf(fi, s)
}

// Reads new container information into config after an update
func updateConfig() {
	log.Println("Updating the config...")
	list := listContainers()
	up := []*docker.Container{}
	for _, cfgCon := range clusterConfig {
		for _, con := range list {
			if cfgCon.Name == con.Name {
				up = append(up, con)
			}
		}
	}
	y, err := yaml.Marshal(up)
	if err != nil {
		logErr(err)
	}
	s := string(y)
	fi, err := os.Create(*configFilePath)
	if err != nil {
		logErr(err)
	}
	fmt.Fprintf(fi, s)
}
