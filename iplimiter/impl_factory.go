package iplimiter

import (
	"github.com/derlaft/ratecounter/counter"
	"github.com/derlaft/ratecounter/iface"
	"time"
)

type limiterFactoryImpl struct {
	// Input settings of rate-limiting
	Interval       time.Duration
	Accuracy       time.Duration
	MaxNumberPerIp int
	// Ability to create counters
	CounterFactory iface.CounterFactory
}

func GetFactory(interval, accuracy time.Duration, maxPerInt int) LimiterFactory {
	return &limiterFactoryImpl{
		Interval:       interval,
		Accuracy:       accuracy,
		MaxNumberPerIp: maxPerInt,
		CounterFactory: counter.GetFactory(),
	}
}

func (lf *limiterFactoryImpl) New() Limiter {

	return &limiterImpl{
		Interval:       lf.Interval,
		Accuracy:       lf.Accuracy,
		MaxNumberPerIp: lf.MaxNumberPerIp,
		CounterFactory: lf.CounterFactory,
		Counters:       make(map[string]iface.Counter),
		GlobalCounter:  lf.CounterFactory.New(lf.Interval, lf.Accuracy),
	}
}

func (lf *limiterFactoryImpl) Restore(data []byte) (Limiter, error) {

	l := &limiterImpl{
		Interval:       lf.Interval,
		Accuracy:       lf.Accuracy,
		MaxNumberPerIp: lf.MaxNumberPerIp,
		CounterFactory: lf.CounterFactory,
		Counters:       make(map[string]iface.Counter),
	}

	err := l.LoadState(data)
	if err != nil {
		return nil, err
	}

	return l, nil
}
