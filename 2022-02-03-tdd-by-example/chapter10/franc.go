package chapter10

type Franc struct {
	*money
}

func NewFranc(a int, c string) *Franc {
	return &Franc{
		NewMoney(a, c),
	}
}
