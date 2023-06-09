package coinflip

import "time"

type Clock interface {
	Now() time.Time
}

func WithClock(cl Clock) func(*Coinflip) {
	return func(cf *Coinflip) {
		cf.clock = cl
	}
}

type ManualClock struct {
	time time.Time
}

func NewManualClock() *ManualClock {
	return &ManualClock{time: time.Now()}
}

func (c *ManualClock) Advance(d time.Duration) {
	c.time = c.time.Add(d)
}

func (c *ManualClock) Now() time.Time {
	return c.time
}

type RegularClock struct{}

func (c *RegularClock) Now() time.Time {
	return time.Now()
}
