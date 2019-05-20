package main

import (
	"fmt"
	"github.com/derlaft/ratecounter/iplimiter"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

const (
	// filename where counter data will be saved
	filename   = "state.rtt"
	fileMode   = 0644
	listenAddr = "0.0.0.0:8081"
	windowSize = time.Second * 20
	// how accurate the counter should be
	// every request time is rounded to this precision
	accuracy = time.Second / 5
	// request cut-out limit
	maxRequests = 15
)

func main() {

	factory := iplimiter.GetFactory(windowSize, accuracy, maxRequests)

	var i iplimiter.Limiter

	data, err := ioutil.ReadFile(filename)
	if os.IsNotExist(err) {
		// file not found
		// create an empty counter
		i = factory.New()
	} else if err != nil {
		// any other error
		// probably not a good idea to die here
		// depends on usage pattern
		log.Fatal("Failed to restore state", err)
	} else {
		// data is loaded, we can just restore it
		i, err = factory.Restore(data)
		if err != nil {
			log.Fatal("Failed to restore state", err)
		}
	}

	// this will handle ^C
	go signalHandler(i)

	// this will save data to disk periodically
	go checkpointSave(i)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		totalRequests, err := processUserRequest(i, w, r)
		if err != nil {

			// write a bad respone code
			w.WriteHeader(http.StatusInternalServerError)

			log.Println("Warning: error while processing request:", err)

			_, err = fmt.Fprintf(w, fmt.Sprintf("Error while processing your request: %v\n", err))
			if err != nil {
				log.Println("Warning: http response write error:", err)
			}

			return
		}

		_, err = fmt.Fprintf(w, fmt.Sprintf("%05d\n", totalRequests))
		if err != nil {
			log.Println("Warning: http response write error:", err)
		}

	})

	log.Println("Started listening on", listenAddr)
	err = http.ListenAndServe(listenAddr, nil)
	if err != nil {
		log.Fatal(err)
	}

}

func processUserRequest(i iplimiter.Limiter, w http.ResponseWriter, r *http.Request) (int, error) {

	ip, err := determineUserIP(r)
	if err != nil {
		return 0, fmt.Errorf("Warning: error while determining the user IP: %v", err)
	}

	// parse user IP
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return 0, fmt.Errorf("Warning: error while parsing IP %v", r.RemoteAddr)
	}

	// check if we pass this user throught
	shouldReject := i.OnRequest(parsedIP)
	if shouldReject {
		return 0, fmt.Errorf("Too many requests from your IP (%v), go away", ip)
	}

	// answer the client safely
	return i.TotalRequests(), nil
}
