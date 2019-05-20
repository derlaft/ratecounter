package iplimiter

import (
	"net"
)

type Limiter interface {
	// OnRequest must be called after each request. if bool=true is returned, the request must be rejected
	OnRequest(ip net.IP) (limitReached bool)

	// Returns total number of requests
	TotalRequests() int

	// Cleanup removes old unused data, must be run periodically
	Cleanup()

	// Marshal internal state for saving somewhere
	SaveState() ([]byte, error)
}

type LimiterFactory interface {

	// Creates a new counter
	New() Limiter

	// Restores a previously-saved counter from a file
	Restore([]byte) (Limiter, error)
}
