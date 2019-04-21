package counter

import (
	"time"
)

// item contains an info about number of increments in each c.Accuracy interval
type item struct {
	// Sequential number of interval (since ~~moon landing)
	Timeint int64
	// Number of increments here
	Count int
}

// now returns sequentional number of current time interval
func (c *counter) now() int64 {
	return time.Now().UnixNano() / c.Accuracy.Nanoseconds()
}

// window returns number of valid meausrement units
func (c *counter) window() int64 {
	return c.WindowSize.Nanoseconds() / c.Accuracy.Nanoseconds()
}
