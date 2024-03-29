package chapter11

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiplication(t *testing.T) {
	five := dollar(5)
	assert.True(t, dollar(10).Equals(five.Times(2))) // TODO replace with equals because there is no struct equality self-definition
	assert.True(t, dollar(15).Equals(five.Times(3)))
}

func TestEquality(t *testing.T) {
	assert.True(t, dollar(5).Equals(dollar(5)))
	assert.False(t, dollar(5).Equals(dollar(6)))
	assert.False(t, franc(5).Equals(dollar(5)))
}
