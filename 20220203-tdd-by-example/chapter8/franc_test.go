package chapter8

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFrancMultiplication(t *testing.T) {
	five := franc(5)
	assert.Equal(t, franc(10), five.Times(2))
	assert.Equal(t, franc(15), five.Times(3))
}
