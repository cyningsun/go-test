package chapter16

type Sum struct {
	Expression
	augend Expression
	added  Expression
}

func NewSum(augend, added Expression) *Sum {
	return &Sum{
		augend: augend,
		added:  added,
	}
}

func (s *Sum) Reduce(bank *Bank, to string) Money {
	amount := s.augend.Reduce(bank, to).Amount() + s.added.Reduce(bank, to).Amount()
	return NewMoney(amount, to)
}

func (s *Sum) Plus(added Expression) Expression {
	return NewSum(s, added)
}

func (s *Sum) Times(multiplier int) Expression {
	return NewSum(s.augend.Times(multiplier), s.added.Times(multiplier))
}
