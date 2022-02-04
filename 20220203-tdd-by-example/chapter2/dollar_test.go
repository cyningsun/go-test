package chapter2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiplication(t *testing.T) {
	five := NewDollar(5)
	product := five.Times(2)
	assert.Equal(t, 10, product.amount)
	product = five.Times(3)
	assert.Equal(t, 15, product.amount)
}
