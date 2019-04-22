package counter

import (
	"errors"
	"github.com/derlaft/ratecounter/iface"
	"github.com/paulbellamy/ratecounter"
	"time"
)

type counter struct {
	C *ratecounter.RateCounter
}

// NewCounter returns a new counter instance with given parameters
func NewCounter(window time.Duration) (iface.Counter, error) {
	return &counter{
		C: ratecounter.NewRateCounter(window),
	}, nil
}

func (c *counter) Incr() {
	c.C.Incr(1)
}

func (c *counter) Count() int {
	return int(c.C.Rate())
}

func (c *counter) Save(string) error {
	return errors.New("Not implemented")
}

func (c *counter) Load(string) error {
	return errors.New("Not implemented")
}
