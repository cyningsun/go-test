package chapter7

type Franc struct {
	Money
}

func NewFranc(a int) *Franc {
	return &Franc{
		NewMoney(a),
	}
}

func (d *Franc) Times(multiplier int) *Franc {
	return NewFranc(d.Amount() * multiplier)
}
