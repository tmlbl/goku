package main

import (
	"encoding/json"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"github.com/ogier/pflag"
	"io/ioutil"
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
	fmt.Println("Parsing the config...")
	data, err := ioutil.ReadFile(*configFilePath)
	if err != nil {
		fmt.Errorf("Couldn't read the config file!", err)
		return false
	}
	json.Unmarshal(data, &clusterConfig)
	printcfg, _ := json.MarshalIndent(clusterConfig, "", "  ")
	fmt.Println(string(printcfg))
	return true
}

// Check the current state against the config state
func checkState(con *docker.Container, cfgCon *docker.Container) {
	up := con.State.Running
	if !up {
		fmt.Println(cfgCon.Name, "is down, bringing it up...")
		bringUpContainer(cfgCon)
	} else {
		fmt.Println(cfgCon.Name, "is up")
	}
}

// Attempts to start a container using its ID
func bringUpContainer(cfgCon *docker.Container) {
	if len(cfgCon.ID) < 10 {
		fmt.Println("Container was not found, creating it...")
		createContainer(cfgCon)
	}
	cli := newClient()
	err := cli.StartContainer(cfgCon.ID, &docker.HostConfig{})
	if err != nil {
		fmt.Println("Error starting the container!")
		fmt.Println(err)
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
		fmt.Println("No containers found! Exiting...")
		os.Exit(1)
	}
	js, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(js))
	fi, err := os.Create(*configFilePath)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(fi, string(js))
}
