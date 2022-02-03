package chapter10

func dollar(amount int) Money {
	return NewDollar(amount, "USD")
}

func franc(amount int) Money {
	return NewFranc(amount, "CHF")
}
