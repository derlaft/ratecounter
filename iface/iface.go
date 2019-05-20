package iface

import (
	"time"
)

//go:generate mockgen -package counter_mocks -destination ../mocks/counter.go github.com/derlaft/ratecounter/iface Counter,CounterFactory

// Counterer allows to keep track of number of events at given time window
type Counter interface {

	// Incr adds a new event
	Incr()

	// Count returns number of events at given time window
	Count() int

	// Save saves data to disk for future use
	Save() ([]byte, error)
}

type CounterFactory interface {

	// Create a new instance of a counter
	New(windowSize, accuracy time.Duration) Counter

	// Load an instance from a saved data
	Load(windowSize, accuracy time.Duration, data []byte) (Counter, error)
}
