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

type ManualClock struct {
	NowValue time.Time
}

func (c *ManualClock) Now() time.Time {
	return c.NowValue
}

type RegularClock struct{}

func (c *RegularClock) Now() time.Time {
	return time.Now()
}

type Coin interface {
	Flip() string
}

func WithCoin(c Coin) func(cf *Coinflip) {
	return func(cf *Coinflip) {
		cf.coin = c
	}
}

type PredictableCoin struct {
	Outcome string
}

func (c *PredictableCoin) Flip() string {
	return c.Outcome
}

type RegularCoin struct{}

func (c *RegularCoin) Flip() string {
	return []string{"head", "tail"}[rand.Intn(2)]
}

type Coinflip struct {
	clock      Clock
	coin       Coin
	isStarted  bool
	ledger     *Ledger
	Outcome    maybe.Maybe[string]
	hasFlipped bool

	bettingDuration time.Duration

	timeLastUpdated time.Time
	timeStarted     time.Time

	MessageChannel chan<- string
}

func NewCoinflip(options ...func(*Coinflip)) *Coinflip {
	c := &Coinflip{
		clock:           &RegularClock{},
		coin:            &RegularCoin{},
		isStarted:       false,
		ledger:          &Ledger{},
		Outcome:         maybe.Nothing[string](),
		hasFlipped:      false,
		bettingDuration: 10 * time.Second,
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
	c.Outcome = maybe.Just(c.coin.Flip())
	c.hasFlipped = true
}

func (c *Coinflip) HasFlipped() bool {
	return c.hasFlipped
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

func (c *Coinflip) Start() {
	c.isStarted = true
	c.timeStarted = c.clock.Now()
}

func (c *Coinflip) Update() {
	now := c.clock.Now()
	if !c.HasFlipped() && c.timeStarted.Add(c.bettingDuration).After(now) {
		c.Flip()
	}
	c.timeLastUpdated = now
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
