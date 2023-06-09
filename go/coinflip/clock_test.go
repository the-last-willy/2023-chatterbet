package coinflip_test

import (
	"chatterbet/coinflip"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_manual_clock_time_can_be_advanced(t *testing.T) {
	someDuration := 13 * time.Hour

	c := coinflip.NewManualClock()

	now := c.Now()
	c.Advance(someDuration)
	then := c.Now()

	assert.Equal(t, someDuration, then.Sub(now))
}
