package counter

import (
	"fmt"
	"github.com/derlaft/ratecounter/iface"
	"sync"
	"time"
)

type counter struct {
	WindowSize time.Duration
	Accuracy   time.Duration
	ProbeVals  []*item
	Lock       sync.RWMutex
	count      int
	relocs     int
}

// NewCounter returns a new counter instance with given parameters
func NewCounter(window, acc time.Duration) (iface.Counter, error) {

	if window == 0 || acc == 0 {
		return nil, fmt.Errorf("Zero duration (window or accuracy)")
	}

	return &counter{
		WindowSize: window,
		Accuracy:   acc,
	}, nil
}

// Incr adds a value to the last item.
// In case ProbeVals limit is reached, old items will be removed
// and data will be copied to a new slice.
func (c *counter) Incr() {

	// under write mutex
	c.Lock.Lock()
	defer c.Lock.Unlock()

	var (
		i      *item
		now    = c.now()
		size   = len(c.ProbeVals)
		lastId = size - 1
	)

	// there are only two cases - last value is current interval and not
	if size == 0 || c.ProbeVals[lastId].Timeint != now {
		// create new element
		i = &item{
			Timeint: now,
			Count:   0,
		}
		c.ProbeVals = append(c.ProbeVals, i)
	} else {

		// just extract element
		i = c.ProbeVals[lastId]
	}

	// increment both local and global counter
	i.Count += 1
	c.count += 1

	// we have reached the limit of this slice
	// a new array will be allocated anyway
	// so better copy values now and get rid of
	if len(c.ProbeVals) == cap(c.ProbeVals) {

		// debugging purposes
		c.relocs += 1

		// old values can cause problems here
		c.cleanup()

		// once relocation per two intervals
		// feels nice
		// but must do math
		var newSize = int(c.window() * 2)

		c.ProbeVals = append(
			make([]*item, 0, newSize),
			c.ProbeVals...,
		)

	}
}

// cleanup removes timed-out items from ProbeVals
// must be used it only under the write lock
// has linear complexity, but rarely does any work at all
func (c *counter) cleanup() {

	var (
		iter        = 0
		now         = c.now()
		cleanBefore = now - c.window()
	)

	// Remove all the useless intervals
	// Decrement global counter
	// This is linear, but should not really happen that often

	for iter < len(c.ProbeVals) && c.ProbeVals[iter].Timeint < cleanBefore {

		// decrement global counter
		c.count -= c.ProbeVals[0].Count

		// remove item
		c.ProbeVals = c.ProbeVals[1:]

		// next iteration
		iter += 1
	}

}

// Count returns a sum of active items
// the number is pre-calculated, but old timed-out values must be removed
// before each call
func (c *counter) Count() int {

	// under read mutex
	c.Lock.Lock()
	defer c.Lock.Unlock()

	// cleanup
	c.cleanup()

	// count is now valid
	return c.count
}
