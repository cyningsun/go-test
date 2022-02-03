package chapter8

type Dollar struct {
	*money
}

func NewDollar(a int) *Dollar {
	d := &Dollar{}
	d.money = NewMoney(d, a)
	return d
}

func (d *Dollar) Times(multiplier int) Money {
	return NewDollar(d.Amount() * multiplier)
}
