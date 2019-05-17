package main

import (
	"errors"
	"github.com/derlaft/ratecounter/iplimiter"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const checkpointInterval = time.Second / 2

// checkpointSave dumps counter data to disk (every checkpointInterval)
func checkpointSave(i iplimiter.Limiter) {

	// @TODO: gracefull termination

	for range time.Tick(checkpointInterval) {

		err := saveToFile(i)
		if err != nil {
			log.Println("Warning: failed to save file", err)
		}

	}
}

// signalHandler terminates server on kill signal
func signalHandler(i iplimiter.Limiter) {

	// bind && wait signal
	var sigs = make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	// save to file
	err := saveToFile(i)
	if err != nil {
		log.Fatal("Failed to save file", err)
	}

	// exit
	// it is probably a good idea to stop http server (and everything else) first when using it on a  real application

	log.Println("Terminating")
	os.Exit(0)
}

func saveToFile(i iplimiter.Limiter) error {
	// cleanup
	i.Cleanup()

	// get data
	data, err := i.SaveState()
	if err != nil {
		return err
	}

	// save it to file
	return ioutil.WriteFile(filename, data, fileMode)
}

// some copy&paste frm production-grade (lolsad) code
func determineUserIP(r *http.Request) (string, error) {

	var ip string

	// trying to parse by remote addr
	remoteTokens := strings.Split(r.RemoteAddr, ":")
	// ipv6 problems, yay
	if len(remoteTokens) > 0 && !strings.HasPrefix(remoteTokens[0], "[") {
		ip = remoteTokens[0]
	}

	// get ip by x-forwarded-for
	if forwarded := r.Header.Get("X-Real-Ip"); forwarded > "" {
		ip = forwarded
	}

	// if IP is empty, we could not determine the source IP; fallback to default
	if ip == "" {
		return "", errors.New("failed to determine source IP address")
	}

	return ip, nil
}
