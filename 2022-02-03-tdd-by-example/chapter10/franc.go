package chapter10

type Franc struct {
	*money
}

func NewFranc(a int, c string) *Franc {
	return &Franc{
		NewMoney(a, c),
	}
}

// func (f *Franc) Times(multiplier int) Money {
// 	return NewFranc(f.Amount()*multiplier, f.currency)
// }
