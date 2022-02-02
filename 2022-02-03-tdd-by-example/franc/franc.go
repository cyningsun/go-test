package franc

import (
	"github.com/cyningsun/go-test/2022-02-03-tdd-by-example/money"
)

type Franc struct {
	*money.Money
}

func NewFranc(a int) *Franc {
	return &Franc{
		money.NewMoney(a),
	}
}

func (f *Franc) times(multiplier int) *Franc {
	return NewFranc(f.Amount() * multiplier)
}
