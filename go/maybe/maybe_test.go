package maybe_test

import (
	"chatterbet/maybe"

	"github.com/stretchr/testify/assert"

	"testing"
)

func TestMaybeCanBeNothing(t *testing.T) {
	m := maybe.Nothing[any]()
	_, h := m.Value()
	assert.False(t, h)
}

func TestMaybeCanBeJustAValue(t *testing.T) {
	m := maybe.Just(7)
	v, has := m.Value()
	assert.True(t, has)
	assert.Equal(t, v, 7)
}
