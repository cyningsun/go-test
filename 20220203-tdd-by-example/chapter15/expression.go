package chapter14

type Expression interface {
	Reduce(bank *Bank, to string) Money
	Plus(Expression) Expression
}
