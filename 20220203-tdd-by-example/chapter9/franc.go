package chapter9

type Franc struct {
	*money
}

func NewFranc(a int, c string) *Franc {
	f := &Franc{}
	f.money = NewMoney(f, a, c)
	return f
}

func (f *Franc) Times(multiplier int) Money {
	return franc(f.Amount() * multiplier)
}
