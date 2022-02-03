package chapter8

type Money interface {
	Amount() int
	Equals(interface{}) bool
	Times(multiplier int) Money
}

type money struct {
	amount int
}

func NewMoney(a int) *money {
	return &money{
		amount: a,
	}
}

func (m *money) Amount() int {
	return m.amount
}

func (m *money) Equals(obj interface{}) bool {
	sec := obj.(Money)
	return m.Amount() == sec.Amount()
}

func (m *money) Times(multiplier int) Money {
	return NewDollar(m.Amount() * multiplier)
}

func dollar(amount int) Money {
	return NewDollar(amount)
}

func franc(amount int) Money {
	return NewFranc(amount)
}
