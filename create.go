package main

import (
	//"encoding/json"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	//"github.com/ogier/pflag"
	//"io/ioutil"
	//"os"
)

func createContainer(con *docker.Container) {
	opts := docker.CreateContainerOptions{}
	opts.Name = con.Name
	opts.Config = &docker.Config{}
	opts.Config.Cmd = con.Args
	opts.Config.PortSpecs = translatePorts(con)
	opts.Config.Volumes = translateVolumes(con)
	cli := newClient()
	con, err := cli.CreateContainer(opts)
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully created container", con)
}

func translatePorts(con *docker.Container) []string {
	ports := []string{}
	for port, prtz := range con.HostConfig.PortBindings {
		ports = append(ports, port.Port()+":"+prtz[0].HostPort)
	}
	fmt.Println("PORTS I GOT:", ports)
	return ports
}

func translateVolumes(con *docker.Container) map[string]struct{} {
	volumes := make(map[string]struct{})
	for hostdir, condir := range con.Volumes {
		key := hostdir + ":" + condir
		volumes[key] = struct{}{}
	}
	return volumes
}
