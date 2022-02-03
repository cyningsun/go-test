package chapter13

type Expression interface {
	reduce(to string) Money
}
