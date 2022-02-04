package chapter16

type Expression interface {
	Reduce(bank *Bank, to string) Money
	Plus(Expression) Expression
	Times(multiplier int) Expression
}
