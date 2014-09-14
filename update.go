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
		pullImg(img)
		for _, con := range list {
			if con.Config.Image == img {
				log.Println("Updating container", con.Name)
				err := cli.StopContainer(con.ID, 0)
				if err != nil {
					logErr(err)
				}
				opts := docker.RemoveContainerOptions{
					ID: con.ID,
				}
				err = cli.RemoveContainer(opts)
				if err != nil {
					logErr(err)
				}
				found = true
				time.Sleep(*updateInterval)
			}
		}
		if !found {
			log.Println("No containers found for image", img)
		}
	}
}

func pullImg(img string) {
	log.Println("Pulling image", img, "...")
	cli := newClient()
	opts := docker.PullImageOptions{
		OutputStream: os.Stdout,
		Repository:   img,
		Tag:          "latest",
	}
	err := cli.PullImage(opts, docker.AuthConfiguration{})
	if err != nil {
		logErr(err)
	}
}
