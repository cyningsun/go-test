package chapter6

type Dollar struct {
	Money
}

func NewDollar(a int) *Dollar {
	return &Dollar{
		NewMoney(a),
	}
}

func (d *Dollar) Times(multiplier int) *Dollar {
	return NewDollar(d.Amount() * multiplier)
}
