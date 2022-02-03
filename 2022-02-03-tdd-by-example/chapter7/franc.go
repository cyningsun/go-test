package chapter7

type Franc struct {
	*money
}

func NewFranc(a int) *Franc {
	f := &Franc{}
	f.money = NewMoney(f, a)
	return f
}

func (f *Franc) Times(multiplier int) *Franc {
	return NewFranc(f.Amount() * multiplier)
}
