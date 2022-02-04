package chapter14

type Sum struct {
	Expression
	augend Money
	added  Money
}

func NewSum(augend, added Money) *Sum {
	return &Sum{
		augend: augend,
		added:  added,
	}
}

func (s *Sum) Reduce(bank *Bank, to string) Money {
	amount := s.augend.Amount() + s.added.Amount()
	return NewMoney(amount, to)
}
