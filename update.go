// Implements a simple rolling update by killing containers on an interval
package main

import (
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"github.com/ogier/pflag"
	"os"
	"time"
)

var updateInterval = pflag.Duration(
	"update-interval",
	1000*time.Millisecond,
	"The time interval for rolling updates",
)

func rollingUpdate() {
	fmt.Println("Starting a rolling update...")
	cli := newClient()
	list := listContainers()
	for _, name := range os.Args[2:] {
		time.Sleep(*updateInterval)
		found := false
		for _, con := range list {
			if con.Name[1:] == name {
				fmt.Println("Updating", name)
				err := cli.StopContainer(con.ID, 0)
				if err != nil {
					panic(err)
				}
				opts := docker.RemoveContainerOptions{
					ID: con.ID,
				}
				err = cli.RemoveContainer(opts)
				if err != nil {
					panic(err)
				}
				found = true
			}
		}
		if !found {
			fmt.Println(name, ": Container not found!")
		}
	}
	os.Exit(0)
}
