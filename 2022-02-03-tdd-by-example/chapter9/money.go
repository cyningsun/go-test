package chapter9

import (
	"reflect"
)

type Money interface {
	Amount() int
	Equals(interface{}) bool
	Times(multiplier int) Money
	Currency() string
}

type money struct {
	parent   Money // TODO any way get parent type ?
	amount   int
	currency string
}

func NewMoney(m Money, a int, c string) *money {
	return &money{
		parent:   m,
		amount:   a,
		currency: c,
	}
}

func (m *money) Amount() int {
	return m.amount
}

func (m *money) Equals(obj interface{}) bool {
	sec := obj.(Money)
	return m.Amount() == sec.Amount() && (reflect.ValueOf(obj).Type() == reflect.TypeOf(m.parent))
}

func (m *money) Currency() string {
	return m.currency
}
