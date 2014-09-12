// Implements a simple rolling update by killing containers on an interval
package main

import (
	"github.com/fsouza/go-dockerclient"
	"github.com/ogier/pflag"
	"log"
	"os"
	"time"
)

var updateInterval = pflag.Duration(
	"update-interval",
	5000*time.Millisecond,
	"The time interval for rolling updates",
)

func rollingUpdate(imgs []string) {
	log.Println("Starting a rolling update...")
	cli := newClient()
	list := listContainers()
	for _, img := range imgs {
		found := false
		for _, con := range list {
			if con.Config.Image == img {
				log.Println("Updating container", con.Name)
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
				time.Sleep(*updateInterval)
			}
		}
		if !found {
			log.Println("No containers found for image", img)
		}
	}
	os.Exit(0)
}
