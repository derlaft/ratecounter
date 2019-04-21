package nanocounter

import (
	"time"
)

// Time is stored as nanosecond unixtime
// There's a good reason not to use time.Time:
// - It's 2-3 times bigger
// And while the major work of this module is allocating mem
// it's a good idea to allocate as low mem as possible
// As a result, everything works much faster than with time.Time
// I've also tried to use float64 seconds, but it has terrible precision
type timestamp int64

func Now() timestamp {
	return timestamp(time.Now().UnixNano())
}

func (t timestamp) After(another timestamp) bool {
	return t > another
}

func (t timestamp) Sub(d time.Duration) timestamp {
	return t - timestamp(d.Nanoseconds())
}
