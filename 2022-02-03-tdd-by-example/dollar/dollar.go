package dollar

import (
	"github.com/cyningsun/go-test/2022-02-03-tdd-by-example/money"
)

type Dollar struct {
	*money.Money
}

func NewDollar(a int) *Dollar {
	return &Dollar{
		money.NewMoney(a),
	}
}

func (d *Dollar) times(multiplier int) *Dollar {
	return NewDollar(d.Amount() * multiplier)
}
