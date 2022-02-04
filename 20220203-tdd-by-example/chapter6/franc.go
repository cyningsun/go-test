package chapter6

type Franc struct {
	*money
}

func NewFranc(a int) *Franc {
	return &Franc{
		NewMoney(a),
	}
}

func (f *Franc) Times(multiplier int) *Franc {
	return NewFranc(f.Amount() * multiplier)
}
