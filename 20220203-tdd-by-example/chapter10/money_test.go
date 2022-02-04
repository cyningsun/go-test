package chapter10

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDifferentClassEquality(t *testing.T) {
	assert.True(t, NewMoney(10, "CHF").Equals(NewFranc(10, "CHF")))
}
