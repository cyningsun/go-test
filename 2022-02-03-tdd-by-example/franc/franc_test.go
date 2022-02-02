package franc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiplication(t *testing.T) {
	five := NewFranc(5)
	assert.Equal(t, NewFranc(10), five.times(2))
	assert.Equal(t, NewFranc(15), five.times(3))
}

func TestEquality(t *testing.T) {
	assert.True(t, NewFranc(5).Equals(NewFranc(5)))
}
