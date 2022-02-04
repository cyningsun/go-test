package chapter12

type Money interface {
	Amount() int
	Equals(interface{}) bool
	Times(multiplier int) Money
	Currency() string
	Plus(Money) Expression
	Expression
}

type money struct {
	amount   int
	currency string
}

func NewMoney(a int, c string) *money {
	return &money{
		amount:   a,
		currency: c,
	}
}

func (m *money) Amount() int {
	return m.amount
}

func (m *money) Equals(obj interface{}) bool {
	sec := obj.(Money)
	return m.Amount() == sec.Amount() && m.currency == sec.Currency()
}

func (m *money) Currency() string {
	return m.currency
}

func (m *money) Times(multiplier int) Money {
	return NewMoney(m.Amount()*multiplier, m.currency)
}

func (m *money) Plus(added Money) Expression {
	return NewMoney(m.Amount()+added.Amount(), m.currency)
}
