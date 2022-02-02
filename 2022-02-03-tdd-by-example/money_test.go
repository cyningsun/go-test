package money

import (
	"testing"

	"github.com/cyningsun/go-test/2022-02-03-tdd-by-example/dollar"
	"github.com/cyningsun/go-test/2022-02-03-tdd-by-example/franc"
	"github.com/stretchr/testify/assert"
)

func TestEquality(t *testing.T) {
	assert.True(t, dollar.NewDollar(5).Equals(dollar.NewDollar(5)))
	assert.False(t, dollar.NewDollar(5).Equals(dollar.NewDollar(6)))
	assert.True(t, franc.NewFranc(5).Equals(franc.NewFranc(5)))
	assert.False(t, franc.NewFranc(5).Equals(franc.NewFranc(6)))
	assert.False(t, dollar.NewDollar(5).Equals(franc.NewFranc(5)))
}
