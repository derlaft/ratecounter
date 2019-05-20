package counter

import (
	"github.com/derlaft/ratecounter/iface"
	"time"
)

// factory is just a wrapper for convinience of testing
type factory struct {
}

func GetFactory() iface.CounterFactory {
	return &factory{}
}

func (f *factory) New(windowSize, accuracy time.Duration) iface.Counter {
	return &counter{
		WindowSize: windowSize,
		Accuracy:   accuracy,
	}
}

func (f *factory) Load(windowSize, accuracy time.Duration, data []byte) (iface.Counter, error) {

	c := &counter{
		WindowSize: windowSize,
		Accuracy:   accuracy,
	}

	err := c.Load(data)
	if err != nil {
		return nil, err
	}

	return c, nil
}
