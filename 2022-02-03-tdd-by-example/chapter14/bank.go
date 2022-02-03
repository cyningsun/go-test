package chapter14

type Bank struct {
	rates map[pair]int
}

func NewBank() *Bank {
	return &Bank{
		rates: make(map[pair]int, 0),
	}
}

func (b *Bank) reduce(source Expression, to string) Money {
	return source.reduce(b, to)
}

func (b *Bank) rate(from, to string) int {
	if from == to {
		return 1
	}
	p := NewPair(from, to)
	rate := b.rates[*p]
	return rate
}

func (b *Bank) addRate(from, to string, rate int) {
	p := NewPair(from, to)
	b.rates[*p] = rate
}
