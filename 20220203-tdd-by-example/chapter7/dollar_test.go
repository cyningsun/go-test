package chapter7

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiplication(t *testing.T) {
	five := NewDollar(5)
	assert.Equal(t, NewDollar(10), five.Times(2))
	assert.Equal(t, NewDollar(15), five.Times(3))
}

func TestEquality(t *testing.T) {
	assert.True(t, NewDollar(5).Equals(NewDollar(5)))
	assert.False(t, NewDollar(5).Equals(NewDollar(6)))

	assert.True(t, NewFranc(5).Equals(NewFranc(5)))
	assert.False(t, NewFranc(5).Equals(NewFranc(6)))

	assert.False(t, NewFranc(5).Equals(NewDollar(5)))
}
