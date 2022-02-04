package chapter16

type Bank struct {
	rates map[pair]int
}

func NewBank() *Bank {
	return &Bank{
		rates: make(map[pair]int, 0),
	}
}

func (b *Bank) Reduce(source Expression, to string) Money {
	return source.Reduce(b, to)
}

func (b *Bank) Rate(from, to string) int {
	if from == to {
		return 1
	}
	p := NewPair(from, to)
	rate := b.rates[*p]
	return rate
}

func (b *Bank) AddRate(from, to string, rate int) {
	p := NewPair(from, to)
	b.rates[*p] = rate
}
