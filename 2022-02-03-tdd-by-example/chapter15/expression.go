package chapter14

type Expression interface {
	reduce(bank *Bank, to string) Money
}
