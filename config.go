package main

import (
	"encoding/json"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"github.com/ogier/pflag"
	"io/ioutil"
	"os"
	"strings"
)

var configFilePath = pflag.String(
	"config",
	"./containers.json",
	"The container list to maintain",
)

type Config []docker.APIContainers

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

// Check the current state agains the config state
func checkState(con docker.APIContainers, cfgCon docker.APIContainers) {
	up := strings.Contains(con.Status, "Up")
	if !up {
		fmt.Println(cfgCon.Names, "is down, bringing it up...")
		bringUpContainer(cfgCon)
	} else {
		fmt.Println(cfgCon.Names, "is up")
	}
}

// Attempts to start a container using its ID
func bringUpContainer(cfgCon docker.APIContainers) {
	cli := newClient()
	err := cli.StartContainer(cfgCon.ID, &docker.HostConfig{})
	if err != nil {
		fmt.Println("Error starting the container!")
		fmt.Println(err)
	}
}

// Get the list of running containers
func listContainers() []docker.APIContainers {
	cli := newClient()
	opts := docker.ListContainersOptions{
		All: true,
	}
	list, err := cli.ListContainers(opts)
	if err != nil {
		panic(err)
	}
	return list
}

// Create a new config file based on the current state if none is given
func createConfig() {
	list := listContainers()
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
