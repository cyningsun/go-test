package chapter6

type Money interface {
	Amount() int
	Equals(interface{}) bool
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
