package money

type Amount interface {
	Amount() int
}

type Money struct {
	amount int
}

func NewMoney(a int) *Money {
	return &Money{
		amount: a,
	}
}

func (m *Money) Amount() int {
	return m.amount
}

func (m *Money) Equals(a Amount) bool {
	return m.Amount() == a.Amount()
}
