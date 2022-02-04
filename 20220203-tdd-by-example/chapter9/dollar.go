package chapter9

type Dollar struct {
	*money
}

func NewDollar(a int, c string) *Dollar {
	d := &Dollar{}
	d.money = NewMoney(d, a, c)
	return d
}

func (d *Dollar) Times(multiplier int) Money {
	return dollar(d.Amount() * multiplier)
}
