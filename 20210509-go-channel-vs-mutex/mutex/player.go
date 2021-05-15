package mutex

type Player interface {
	NextScore() (score int, err error)
}
