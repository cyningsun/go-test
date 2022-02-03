package chapter10

type Dollar struct {
	*money
}

func NewDollar(a int, c string) *Dollar {
	return &Dollar{
		NewMoney(a, c),
	}
}

func (d *Dollar) Times(multiplier int) Money {
	return NewDollar(d.Amount()*multiplier, d.currency)
}
