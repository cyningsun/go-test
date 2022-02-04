package chapter4

type Dollar struct {
	amount int
}

func NewDollar(a int) *Dollar {
	return &Dollar{
		amount: a,
	}
}

func (d *Dollar) Times(multiplier int) *Dollar {
	return NewDollar(d.amount * multiplier)
}

func (d *Dollar) Equals(obj interface{}) bool {
	dollar := obj.(*Dollar)
	return d.amount == dollar.amount
}
