package chapter8

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiplication(t *testing.T) {
	five := dollar(5)
	assert.Equal(t, dollar(10), five.Times(2))
	assert.Equal(t, dollar(15), five.Times(3))
}

func TestEquality(t *testing.T) {
	assert.True(t, dollar(5).Equals(dollar(5)))
	assert.False(t, dollar(5).Equals(dollar(6)))

	assert.True(t, franc(5).Equals(franc(5)))
	assert.False(t, franc(5).Equals(franc(6)))

	assert.False(t, franc(5).Equals(dollar(5)))
}
