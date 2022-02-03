package chapter8

type Dollar struct {
	Money
}

func NewDollar(a int) *Dollar {
	return &Dollar{
		NewMoney(a),
	}
}

func (d *Dollar) Times(multiplier int) Money {
	return NewDollar(d.Amount() * multiplier)
}
