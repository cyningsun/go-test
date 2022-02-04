package chapter10

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFrancMultiplication(t *testing.T) {
	five := franc(5)
	assert.True(t, franc(10).Equals(five.Times(2)))
	assert.True(t, franc(15).Equals(five.Times(3)))
}
