package chapter5

type Franc struct {
	amount int
}

func NewFranc(a int) *Franc {
	return &Franc{
		amount: a,
	}
}

func (d *Franc) Times(multiplier int) *Franc {
	return NewFranc(d.amount * multiplier)
}

func (d *Franc) Equals(obj interface{}) bool {
	franc := obj.(*Franc)
	return d.amount == franc.amount
}
