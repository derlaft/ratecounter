package main

import (
	"github.com/derlaft/ratecounter/iface"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const checkpointInterval = time.Second / 2

// checkpointSave dumps counter data to disk (every checkpointInterval)
func checkpointSave(i iface.Counter) {

	// @TODO: gracefull termination

	for range time.Tick(checkpointInterval) {

		err := i.Save(filename)
		if err != nil {
			log.Println("Warning: failed to save file", err)
		}

	}
}

// signalHandler terminates server on kill signal
func signalHandler(i iface.Counter) {

	// bind && wait signal
	var sigs = make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	// save to file
	err := i.Save(filename)
	if err != nil {
		log.Fatal("Failed to save file", err)
	}

	// exit
	// it is probably a good idea to stop http server (and everything else) first when using it on a  real application

	log.Println("Terminating")
	os.Exit(0)
}
