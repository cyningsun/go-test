package chapter16

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

func TestSimpleAddition(t *testing.T) {
	sum := dollar(5).Plus(dollar(5))
	assert.Equal(t, dollar(10), sum)

	bank := NewBank()
	reduced := bank.Reduce(sum, "USD")
	assert.Equal(t, dollar(10), reduced)
}

func TestPlusReturnSum(t *testing.T) {
	five := dollar(5)
	result := five.Plus(five)

	sum := result.(*Sum)
	assert.Equal(t, five, sum.augend)
	assert.Equal(t, five, sum.added)
}

func TestReduceSum(t *testing.T) {
	sum := NewSum(dollar(3), dollar(4))
	bank := NewBank()
	result := bank.Reduce(sum, "USD")
	assert.Equal(t, dollar(7), result)
}

func TestReduceMoney(t *testing.T) {
	bank := NewBank()
	result := bank.Reduce(dollar(1), "USD")
	assert.Equal(t, dollar(1), result)
}

func TestReduceMoneyDifferenceCurrency(t *testing.T) {
	bank := NewBank()
	bank.AddRate("CHF", "USD", 2)
	result := bank.Reduce(franc(2), "USD")

	assert.Equal(t, dollar(1), result)
}

func TestIdentifyRate(t *testing.T) {
	assert.Equal(t, 1, NewBank().Rate("USD", "USD"))
}

func TestMixedAddition(t *testing.T) {
	var fiveBucks Expression = dollar(5)
	var tenFrancs Expression = franc(10)

	bank := NewBank()
	bank.AddRate("CHF", "USD", 2)

	result := bank.Reduce(fiveBucks.Plus(tenFrancs), "USD")
	assert.Equal(t, dollar(10), result)
}

func TestSumPlusMoney(t *testing.T) {
	var fiveBucks Expression = dollar(5)
	var tenFrancs Expression = franc(10)
	var bank *Bank = NewBank()
	bank.AddRate("CHF", "USD", 2)
	var sum Expression = NewSum(fiveBucks, tenFrancs).Plus(fiveBucks)
	var result Money = bank.Reduce(sum, "USD")
	assert.Equal(t, dollar(15), result)
}

func TestSumTimes(t *testing.T) {
	var fiveBucks Expression = dollar(5)
	var tenFrancs Expression = franc(10)
	var bank *Bank = NewBank()
	bank.AddRate("CHF", "USD", 2)
	var sum Expression = NewSum(fiveBucks, tenFrancs).Times(2)
	var result Money = bank.Reduce(sum, "USD")
	assert.Equal(t, dollar(20), result)
}

// func TestPlusSameCurrencyReturnsMoney(t *testing.T) {
// 	var sum Expression = dollar(1).Plus(dollar(1))
// 	_, ok := sum.(*money)
// 	assert.True(t, ok)
// }
