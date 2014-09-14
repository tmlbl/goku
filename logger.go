package main

import (
	"github.com/mgutz/ansi"
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
	log.Println("Initializing logger...")
	if *logFilePath == "" {
		log.Println("No logfile given, will use stdout")
		log.SetOutput(os.Stdout)
	} else {
		log.Println("Using log file")
		fi, err := os.Create(*logFilePath)
		if err != nil {
			logErr(err)
		}
		log.SetOutput(fi)
	}
}

func logErr(err error) {
	se := err.Error()
	log.Println(ansi.Color("Received an error:", "red"), "\n\t", se)
}
