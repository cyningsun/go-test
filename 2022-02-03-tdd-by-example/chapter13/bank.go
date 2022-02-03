package chapter13

type Bank struct {
}

func NewBank() *Bank {
	return &Bank{}
}

func (b *Bank) reduce(source Expression, to string) Money {
	return source.reduce(to)
}
