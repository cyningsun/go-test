package chapter13

func dollar(amount int) Money {
	return NewMoney(amount, "USD")
}

func franc(amount int) Money {
	return NewMoney(amount, "CHF")
}
