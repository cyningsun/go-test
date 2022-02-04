package chapter8

func dollar(amount int) Money {
	return NewDollar(amount)
}

func franc(amount int) Money {
	return NewFranc(amount)
}
