shgod
=====

A lightweight and simple tool that facilitates continuous deployment of docker containers.

## Quick Start
### Installation
Installation is as simple as downloading the binary and making it executable. On a Debian system, you might run something like this:
````
sudo curl -o /usr/bin/shgod https://s3-us-west-2.amazonaws.com/shgod/shgod \
  && sudo chmod +x /usr/bin/shgod
````
Alternatively, if you have a go enviroment on your system, it can be easily compiled from source:
````
go get github.com/tmlbl/shgod &&
cd $GOPATH/src/github.com/tmlbl/shgod &&
go get ./... &&
go build &&
go install
````
The default behavior is to take a "snapshot" of the running containers on your system, and to continously monitor and maintain that state. In an environment where you are running docker containers, run:
````
shgod init
````
This will create a new `containers.yml` file in the current directory, and start the program. shgod runs a heartbeat to check that the specified containers still exist and are still running, and if they are not, it brings the system back to the state described in `containers.yml`. It also runs an http server which can respond to webhooks from the docker hub.

### Rolling Updates

Because shgod will create an identical new container if a container described in `containers.yml` is not found, containers can be removed safely and will be recreated using the corresponding image. In another terminal, run:
````
shgod update <image_name>
````
Where `image_name` is the name of a docker repository. shgod will pull the repository, and then systematically destroy the containers using that image. On the next heartbeat, they are recreated with identical configuration and the latest version of the image.


This same operation can be triggered by a webhook. By default, shgod listens on port `6600` for webhooks. This and many other settings can be configured using flags. Run `shgod --help` for more information.

## Running in the Background
Currently, the easiest way to do this is to simply fork the process:
````
shgod > logfile.log &
````
