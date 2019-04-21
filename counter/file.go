package counter

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const fileMode = 0644

// Save data to disk using json encoding
func (c *counter) Save(filename string) (err error) {

	// locks here (don't want two saves at the same time)
	c.Lock.Lock()
	defer c.Lock.Unlock()

	// it's probably one case we need to drop old values manually
	c.cleanup()

	// want to write at least one value
	if len(c.ProbeVals) <= 0 {
		return nil
	}

	// why json? it's easy to use
	// and this time we don't have much data to write
	bytes, err := json.Marshal(&c.ProbeVals)
	if err != nil {
		return
	}

	return ioutil.WriteFile(filename, bytes, fileMode)
}

// Load json-encoded data from disk
// Clean old values after the load
// Re-calculate sum counter.
func (c *counter) Load(filename string) error {

	// read and decode file
	data, err := ioutil.ReadFile(filename)
	if os.IsNotExist(err) {
		// file not found - exit without error
		return nil
	} else if err != nil {
		return err
	}

	// decode values
	err = json.Unmarshal(data, &c.ProbeVals)
	if err != nil {
		return err
	}

	// it would probably be a good idea to sort && validate values
	// just in case some weird issues (like clock changing between restarts)

	// no need to lock - used only at startup
	// but anyway use it
	c.Lock.Lock()
	defer c.Lock.Unlock()

	// calculate correct common counter
	for _, i := range c.ProbeVals {
		c.count += i.Count
	}

	// cleanup old values
	c.cleanup()

	return nil
}
