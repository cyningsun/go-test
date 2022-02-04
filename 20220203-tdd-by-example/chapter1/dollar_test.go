package chapter1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiplication(t *testing.T) {
	five := NewDollar(5)
	five.Times(2)
	assert.Equal(t, 10, five.amount)
}
