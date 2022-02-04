package chapter7

import (
	"reflect"
)

type Money interface {
	Amount() int
	Equals(interface{}) bool
}

type money struct {
	amount int
	parent Money // TODO any way get parent type ?
}

func NewMoney(m Money, a int) *money {
	return &money{
		parent: m,
		amount: a,
	}
}

func (m *money) Amount() int {
	return m.amount
}

func (m *money) Equals(obj interface{}) bool {
	sec := obj.(Money)
	return m.Amount() == sec.Amount() && (reflect.ValueOf(obj).Type() == reflect.TypeOf(m.parent))
}
