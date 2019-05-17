package iplimiter

import (
	"encoding/json"
	"github.com/derlaft/ratecounter/iface"
	"net"
	"sync"
	"time"
)

type limiterImpl struct {
	// Ability to create counters
	CounterFactory iface.CounterFactory
	// Input settings of rate-limiting
	Interval       time.Duration
	Accuracy       time.Duration
	MaxNumberPerIp int
	// Actual per-ip counters
	// Must be cleaned periodically
	CountersLock sync.RWMutex
	Counters     map[string]iface.Counter
	// Global request counter
	GlobalCounter iface.Counter
}

type encodedCounter struct {
	Global json.RawMessage
	Values map[string]json.RawMessage
}

func (c *limiterImpl) OnRequest(ip net.IP) (limitReached bool) {

	var (
		counter    iface.Counter
		counterKey = ip.String()
	)

	// Try finding already saved counter; if it does not exist - create a new one
	c.CountersLock.Lock()

	presentCounter, found := c.Counters[counterKey]
	if found {
		counter = presentCounter
	} else {
		counter = c.CounterFactory.New(c.Interval, c.Accuracy)
		c.Counters[counterKey] = counter
	}

	c.CountersLock.Unlock()

	// increment a counter
	counter.Incr()

	// check if limit is reached
	limitReached = counter.Count() > c.MaxNumberPerIp

	// increment a global counter if not
	if !limitReached {
		c.GlobalCounter.Incr()
	}

	return
}

func (c *limiterImpl) TotalRequests() int {

	return c.GlobalCounter.Count()
}

func (c *limiterImpl) Cleanup() {

	// obtain an rw-lock
	c.CountersLock.Lock()
	defer c.CountersLock.Unlock()

	for key, value := range c.Counters {

		// just delete every value with zero counter
		if value.Count() == 0 {
			delete(c.Counters, key)
		}
	}
}

func (c *limiterImpl) SaveState() ([]byte, error) {

	// obtain ro-lockl
	c.CountersLock.RLock()
	defer c.CountersLock.RUnlock()

	// prepare output
	var result encodedCounter
	result.Values = make(map[string]json.RawMessage)

	// prepare global counter
	encVal, err := c.GlobalCounter.Save()
	if err != nil {
		return nil, err
	}
	result.Global = json.RawMessage(encVal)

	for k, counter := range c.Counters {

		// encode each counter
		encVal, err = counter.Save()
		if err != nil {
			return nil, err
		}

		// save it
		result.Values[k] = json.RawMessage(encVal)
	}

	// encode it
	return json.Marshal(result)
}

func (c *limiterImpl) LoadState(data []byte) error {

	// locks not needed here
	var state encodedCounter
	err := json.Unmarshal(data, &state)
	if err != nil {
		return err
	}

	// load global counter
	globalCounter, err := c.CounterFactory.Load(
		c.Interval,
		c.Accuracy,
		state.Global,
	)
	if err != nil {
		return err
	}

	c.GlobalCounter = globalCounter

	for key, value := range state.Values {

		// import every counter data
		newCounter, err := c.CounterFactory.Load(
			c.Interval,
			c.Accuracy,
			value,
		)
		if err != nil {
			return err
		}

		// &&save it
		c.Counters[key] = newCounter
	}

	return nil
}
