// If a container in the config file has been removed, create a new one
package main

import (
	//"fmt"
	"github.com/fsouza/go-dockerclient"
	"log"
	"strconv"
	"strings"
)

func createContainer(con *docker.Container) {
	newConfig := &docker.Config{}
	newConfig = con.Config
	latest := latestTag(con.Config.Image)
	if latest != "" {
		newConfig.Image = latest
	}
	log.Println("Latest tag is", newConfig.Image)
	// Delete a conflicting container if it exists
	removeContainer(con)
	opts := docker.CreateContainerOptions{
		Name:   con.Name,
		Config: newConfig,
	}
	log.Println("Image in options is", opts.Config.Image)
	cli := newClient()
	newCon, err := cli.CreateContainer(opts)
	if err != nil {
		logErr(err)
	}
	log.Println("Successfully created container for", newCon.Name, "with image", newCon.Image)
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

// Finds the highest numerical tag for the given repo
// Returns the image in the format me/myimg:mytag
func latestTag(repo string) string {
	log.Println("Got image name", repo, "removing tags")
	if strings.Contains(repo, ":") {
		repo = strings.Split(repo, ":")[0]
	}
	cli := newClient()
	imgs, err := cli.ListImages(true)
	if err != nil {
		panic(err)
	}
	highTag := ""
	highNum := float64(0)
	log.Println("Finding latest tag for", repo)
	for _, img := range imgs {
		for _, tag := range img.RepoTags {
			//log.Println("Looking at", tag)
			if strings.Contains(tag, repo) {
				log.Println(tag, "contains", repo)
				tg := strings.Split(tag, ":")[1]
				log.Println("Parsing tag", tg)
				num, err := strconv.ParseFloat(tg, 64)
				if err == nil {
					log.Println("Found tag for", repo, num)
					if num > highNum {
						log.Println("Setting tag to", num)
						highNum = num
						highTag = tag
					}
				}
			}
		}
	}
	return highTag
}
