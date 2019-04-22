package main

import (
	"fmt"
	"github.com/derlaft/ratecounter/counter"
	"log"
	"net/http"
	"time"
)

const (
	// filename where counter data will be saved
	filename   = "state.rtt"
	listenAddr = "127.0.0.1:8081"
	windowSize = time.Second * 60
	// how accurate the counter should be
	// every request time is rounded to this precision
	accuracy = time.Second / 5
)

func main() {

	i, err := counter.NewCounter(windowSize, accuracy)
	if err != nil {
		log.Fatal("Failed to initialize the counter")
	}

	err = i.Load(filename)
	if err != nil {
		// probably not a good idea to die here
		// depends on usage pattern
		log.Fatal("Failed to restore state", err)
	}

	// this will handle ^C
	go signalHandler(i)

	// this will save data to disk periodically
	go checkpointSave(i)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		// I could potentially run this in a separate subroutine
		// however, it does not actually change anything significantly
		// and even makes everything a little bit slower.
		// It would be possible to run a separate coroutine inside
		// and remember the latest item.
		// It's almost free anyway.
		i.Incr()

		// I could also combine i.Incr() and i.Count(), but
		// this would make this module practically useless
		// outside the test task.

		_, err = fmt.Fprintf(w, fmt.Sprintf("%020d", i.Count()))
		if err != nil {
			log.Println("Warning: http response write error:", err)
		}
	})

	log.Println("Started listening")
	err = http.ListenAndServe(listenAddr, nil)
	if err != nil {
		log.Fatal(err)
	}

}
