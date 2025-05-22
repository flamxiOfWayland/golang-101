package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIfFoo(t *testing.T) {
	assert.True(t, IfFoo(1, 1.5, 0.001))
	assert.False(t, IfFoo(1, -0.01, -0.1))
}

func TestSwitchFoo(t *testing.T) {
	assert.True(t, SwtichFoo(1, 1.5, .001))
	assert.False(t, SwtichFoo(1, -0.01, -0.1))
}

func TestCreateSlice(t *testing.T) {
	_, err := CreateSlice(-2)
	assert.Error(t, err)

	s, err := CreateSlice(2)
	assert.Nil(t, err)
	assert.Equal(t, len(s), 2)
	assert.Equal(t, s, []int{0, 0})
}
