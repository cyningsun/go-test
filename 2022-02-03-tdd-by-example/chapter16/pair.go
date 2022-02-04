package chapter16

type pair struct {
	from string
	to   string
}

func NewPair(f, t string) *pair {
	return &pair{
		from: f,
		to:   t,
	}
}

func (p *pair) Equals(sec interface{}) bool {
	s := sec.(*pair)
	return p.from == s.from && p.to == s.to
}
