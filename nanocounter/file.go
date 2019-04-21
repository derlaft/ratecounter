package nanocounter

import (
	"encoding/binary"
	"log"
	"os"
	"sort"
)

// timestampSize is number of bytes of each timestamp
const timestampSize = 8

func (c *counter) Save(filename string) (err error) {

	// locks here (don't want two saves at the same time)
	c.Lock.RLock()
	defer c.Lock.RUnlock()

	// it's probably one case we need to drop old values manually
	var border = c.findLeftBorder()

	// want to write at least one value
	if len(c.ProbeVals)-border <= 0 {
		return nil
	}

	// open file (create if not exits)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}

	// close file gracefully
	defer func() {
		err = file.Close()
		if err != nil {
			return
		}
	}()

	// write only needed values
	err = binary.Write(file, binary.LittleEndian, c.ProbeVals[border:])
	if err != nil {
		return
	}

	return
}

func (c *counter) Load(filename string) error {

	// open file
	file, err := os.OpenFile(filename, os.O_RDONLY, 0755)
	if os.IsNotExist(err) {
		// if file not found - nothing to worry about
		return nil
	} else if err != nil {
		return err
	}

	// defer fileclose
	defer func() {
		err := file.Close()
		if err != nil {
			log.Printf("Warning: error while closing file: %v", err)
		}
	}()

	// get size info
	stat, err := file.Stat()
	if err != nil {
		return err
	}

	// no need to lock - used only at startup
	// but anyway use it
	c.Lock.Lock()
	defer c.Lock.Unlock()

	// binary.Read excepts a determined number of values
	// probably not the best way to do it
	// but the easiest to implement with lowest disk usage
	c.ProbeVals = make([]timestamp, stat.Size()/timestampSize)

	err = binary.Read(file, binary.LittleEndian, &c.ProbeVals)
	if err != nil {
		return err
	}

	// file may be damaged
	// sort values to avoid breaking everything
	sort.Slice(c.ProbeVals, func(i, j int) bool {
		return c.ProbeVals[i] < c.ProbeVals[j]
	})

	return nil
}
