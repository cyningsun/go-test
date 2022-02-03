package chapter10

type Money interface {
	Amount() int
	Equals(interface{}) bool
	Times(multiplier int) Money
	Currency() string
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

func (f *money) Times(multiplier int) Money {
	return NewMoney(f.Amount()*multiplier, f.currency)
}
