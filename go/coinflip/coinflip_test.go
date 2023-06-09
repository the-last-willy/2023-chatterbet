package coinflip_test

import (
	. "chatterbet/coinflip"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"testing"
)

type NewCoinflipTestSuite struct {
	suite.Suite
	cf *Coinflip

	clock *ManualClock
}

func TestNewCoinflipTestSuite(t *testing.T) {
	suite.Run(t, new(NewCoinflipTestSuite))
}

func (suite *NewCoinflipTestSuite) SetupTest() {
	suite.clock = NewManualClock()
	suite.cf = NewCoinflip(WithClock(suite.clock))
}

func (suite *NewCoinflipTestSuite) TestShouldNotBeStarted() {
	assert.False(suite.T(), suite.cf.IsStarted())
}

func (suite *NewCoinflipTestSuite) TestCanBeStarted() {
	suite.cf.Start()
	assert.True(suite.T(), suite.cf.IsStarted())
}

func (suite *NewCoinflipTestSuite) TestShouldStartWhenSomeoneSendPlayMessage() {
	_ = suite.cf.Process(&Message{
		User:    "user#1",
		Content: "!play coinflip",
	})
	assert.True(suite.T(), suite.cf.IsStarted())
}

func (suite *NewCoinflipTestSuite) Test_flips_some_time_after_getting_started() {
	suite.cf.Start()
	suite.clock.Advance(11 * time.Second)
	suite.cf.Update()
	assert.True(suite.T(), suite.cf.HasFlipped())
}

type StartedCoinflipTestSuite struct {
	suite.Suite
	cf *Coinflip
}

func TestStartedCoinflipTestSuite(t *testing.T) {
	suite.Run(t, new(StartedCoinflipTestSuite))
}

func (suite *StartedCoinflipTestSuite) SetupTest() {
	suite.cf = NewCoinflip()
	suite.cf.Start()
}

func (suite *StartedCoinflipTestSuite) TestShouldRegisterABetOnHead() {
	_ = suite.cf.Process(&Message{
		Content: "!bet head",
		User:    "user#12",
	})
	assert.Contains(suite.T(), suite.cf.AllBets(), Bet{
		Outcome: "head",
		User:    "user#12",
	})
}

func (suite *StartedCoinflipTestSuite) TestShouldRegisterABetOnTail() {
	_ = suite.cf.Process(&Message{
		Content: "!bet tail",
		User:    "user#12",
	})
	assert.Contains(suite.T(), suite.cf.AllBets(), Bet{
		Outcome: "tail",
		User:    "user#12",
	})
}

func (suite *StartedCoinflipTestSuite) TestShouldNotHaveAnOutcomeBeforeItsFlipped() {
	_, has := suite.cf.Outcome.Value()
	assert.False(suite.T(), has)
}

func (suite *StartedCoinflipTestSuite) TestShouldHaveAnOutcomeAfterItsFlipped() {
	suite.cf.Flip()
	v, has := suite.cf.Outcome.Value()
	assert.True(suite.T(), has)
	assert.Contains(suite.T(), []string{"head", "tail"}, v)
}

type CoinflipWithSomeBetsTestSuite struct {
	suite.Suite
	cf   *Coinflip
	coin *PredictableCoin
}

func TestCoinflipWithSomeBetsTestSuite(t *testing.T) {
	suite.Run(t, new(CoinflipWithSomeBetsTestSuite))
}

func (suite *CoinflipWithSomeBetsTestSuite) SetupTest() {
	suite.coin = &PredictableCoin{}
	suite.cf = NewCoinflip(WithCoin(suite.coin))
	suite.cf.Start()
}

func (suite *CoinflipWithSomeBetsTestSuite) TestFlippingOnHeadShouldWinBetsOnHead() {
	_ = suite.cf.Process(&Message{
		User:    "user#3",
		Content: "!bet head",
	})
	_ = suite.cf.Process(&Message{
		User:    "user#4",
		Content: "!bet tail",
	})
	suite.coin.Outcome = "head"
	suite.cf.Flip()
	bs := suite.cf.WonBets()
	assert.Contains(suite.T(), bs, Bet{
		Outcome: "head",
		User:    "user#3",
	})
}

func (suite *CoinflipWithSomeBetsTestSuite) TestFlippingOnTailShouldLoseBetsOnHead() {
	_ = suite.cf.Process(&Message{
		User:    "user#3",
		Content: "!bet head",
	})
	_ = suite.cf.Process(&Message{
		User:    "user#4",
		Content: "!bet tail",
	})
	suite.coin.Outcome = "tail"
	suite.cf.Flip()
	bs := suite.cf.LostBets()
	assert.Contains(suite.T(), bs, Bet{
		Outcome: "head",
		User:    "user#3",
	})
}

func Test_coinflip_flips_a_single_time_after_betting_is_over(t *testing.T) {
	cl := &ManualClock{}
	co := &PredictableCoin{Outcome: "head"}
	cf := NewCoinflip(WithClock(cl), WithCoin(co))

	cf.Start()
	cl.Advance(11 * time.Second)
	cf.Update()
	co.Outcome = "tail"
	cf.Update()

	o, _ := cf.Outcome.Value()
	assert.Equal(t, "head", o)
}

func Test_coinflip_sends_a_message_when_it_is_started(t *testing.T) {
	messages := make(chan string, 5)
	cf := NewCoinflip()
	cf.MessageChannel = messages

	cf.Start()
	cf.Update()

	assert.Equal(t, 1, len(messages))
}
