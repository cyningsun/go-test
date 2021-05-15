package csp

type Player interface {
	NextScore() (score int, err error)
}
