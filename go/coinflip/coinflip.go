package coinflip

import (
	"chatterbet/maybe"
	"errors"
	"math/rand"
	"time"
)

type Bet struct {
	Outcome string
	User    string
}

type Clock interface {
	Now() time.Time
}

func WithClock(cl Clock) func(*Coinflip) {
	return func(cf *Coinflip) {
		cf.clock = cl
	}
}

type Coin interface {
	Flip() string
}

type Coinflip struct {
	clock     Clock
	isStarted bool
	ledger    *Ledger
	Outcome   maybe.Maybe[string]
}

func NewCoinflip(options ...func(*Coinflip)) *Coinflip {
	c := &Coinflip{
		clock:     &RegularClock{},
		isStarted: false,
		ledger:    &Ledger{},
	}
	for _, o := range options {
		o(c)
	}
	return c
}

func (c *Coinflip) AllBets() []Bet {
	return c.ledger.Bets
}

func (c *Coinflip) Flip() {
	if rand.Intn(2) == 0 {
		c.Outcome = maybe.Just("head")
	} else {
		c.Outcome = maybe.Just("tail")
	}
}

func (c *Coinflip) LostBets() []Bet {
	out, has := c.Outcome.Value()
	if !has {
		return nil
	} else {
		var bs []Bet
		for _, b := range c.ledger.Bets {
			if b.Outcome != out {
				bs = append(bs, b)
			}
		}
		return bs
	}
}

func (c *Coinflip) WonBets() []Bet {
	out, has := c.Outcome.Value()
	if !has {
		return nil
	} else {
		var bs []Bet
		for _, b := range c.ledger.Bets {
			if b.Outcome == out {
				bs = append(bs, b)
			}
		}
		return bs
	}
}

func (c *Coinflip) IsStarted() bool {
	return c.isStarted
}

func (c *Coinflip) Process(m *Message) error {
	if m.Content == "!bet head" {
		c.ledger.Register(Bet{
			Outcome: "head",
			User:    m.User,
		})
		return nil
	} else if m.Content == "!bet tail" {
		c.ledger.Register(Bet{
			Outcome: "tail",
			User:    m.User,
		})
		return nil
	} else if m.Content == "!play coinflip" {
		c.Start()
		return nil
	} else {
		return errors.New("invalid message")
	}
}

func (c *Coinflip) registerBet(b Bet) {
	c.ledger.Register(b)
}

func (c *Coinflip) Start() {
	c.isStarted = true
}

type Ledger struct {
	Bets []Bet
}

func (l *Ledger) Register(b Bet) {
	l.Bets = append(l.Bets, b)
}

type Message struct {
	User    string
	Content string
}

type RegularClock struct {
}

func (c *RegularClock) Now() time.Time {
	return time.Now()
}
