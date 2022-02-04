package chapter1

type Dollar struct {
	amount int
}

func NewDollar(a int) *Dollar {
	return &Dollar{
		amount: a,
	}
}

func (d *Dollar) Times(multiplier int) {
	d.amount = d.amount * multiplier
}
