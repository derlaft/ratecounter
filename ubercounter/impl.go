package nanocounter

import (
	"fmt"
	"github.com/derlaft/ratecounter/iface"
	"sort"
	"sync"
	"time"
)

type counter struct {
	WindowSize time.Duration
	ProbeVals  []timestamp
	Lock       sync.RWMutex
	Relocs     int
}

func NewCounter(window time.Duration) (iface.Counter, error) {

	if window == 0 {
		return nil, fmt.Errorf("Zero window")
	}

	return &counter{
		WindowSize: window,
	}, nil
}

func (c *counter) Incr() {

	// under write mutex
	c.Lock.Lock()
	defer c.Lock.Unlock()

	c.ProbeVals = append(c.ProbeVals, Now())

	var oldSize = len(c.ProbeVals)

	// we have reached the limit of this slice
	// a new array will be allocated anyway
	// so better copy values now and get rid of
	if oldSize == cap(c.ProbeVals) {

		c.Relocs += 1

		var (
			border = c.findLeftBorder()
			// probably the default append behaviour works better
			// but it's pretty lame to implement
			// and 1,000,000 is already an extreme number of connections per minute
			// so just scale up to 120% of valid values
			// the best case would probably 200% in this case -
			// if expecting the amount of requests to be relatively constant
			// must do math %)
			newSize = (oldSize - border) + (oldSize-border)/2
		)

		c.ProbeVals = append(
			make([]timestamp, 0, newSize),
			c.ProbeVals[border:]...,
		)
	}
}

func (c *counter) findLeftBorder() int {

	var (
		leftBorder = Now().Sub(c.WindowSize)
		numValues  = len(c.ProbeVals)
	)

	// completely no data here
	if numValues == 0 {
		return 0
	}

	var rule = func(i int) bool {
		// Search finds smallest index where func returns true
		// if there's no such elements, it returns
		// a (non-valid) index that is equal to the length

		return c.ProbeVals[i].After(leftBorder)
	}

	return sort.Search(numValues, rule)
}

func (c *counter) Count() int {

	// under read mutex
	c.Lock.RLock()
	defer c.Lock.RUnlock()

	// binary search for the first value
	// should be relatively fast
	return len(c.ProbeVals) - c.findLeftBorder()
}
