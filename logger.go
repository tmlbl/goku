package main

import (
	"github.com/ogier/pflag"
	"log"
	"os"
)

var logFilePath = pflag.String(
	"logfile",
	"",
	"The path to the log file, if any",
)

func initLogger() {
	log.SetFlags(0)
	if *logFilePath == "" {
		log.SetOutput(os.Stdout)
	} else {
		fi, err := os.Create(*logFilePath)
		if err != nil {
			panic(err)
		}
		log.SetOutput(fi)
	}
}
