package iface

// Counterer allows to keep track of number of events at given time window
type Counter interface {

	// Incr adds a new event
	Incr()

	// Count returns number of events at given time window
	Count() int

	// Save saves data to disk for future use
	Save(fname string) error

	// Load restores previously-saved data
	Load(fname string) error
}
