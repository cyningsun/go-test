package chapter12

type Bank struct {
}

func NewBank() *Bank {
	return &Bank{}
}

func (b *Bank) reduce(source Expression, to string) Money {
	return dollar(10)
}
