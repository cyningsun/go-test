package chapter5

type Franc struct {
	amount int
}

func NewFranc(a int) *Franc {
	return &Franc{
		amount: a,
	}
}

func (f *Franc) Times(multiplier int) *Franc {
	return NewFranc(f.amount * multiplier)
}

func (f *Franc) Equals(obj interface{}) bool {
	franc := obj.(*Franc)
	return f.amount == franc.amount
}
