package testcode

import (
	"fmt"
	"math/rand"
)

var (
	Intn = rand.Intn
)

func Hello(name string) string {
	return fmt.Sprintf("hello %v:%v", name, Intn(10))
}
