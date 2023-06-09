package coinflip_test

import (
	"chatterbet/coinflip"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_non_empty_channel_is_empty_after_clear(t *testing.T) {
	c := make(chan int, 1)
	c <- 7

	coinflip.ClearChannel(c)

	assert.Empty(t, c)
}
