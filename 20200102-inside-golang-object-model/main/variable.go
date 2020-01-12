package main

import (
	"fmt"
	"unsafe"
)

var (
	X, Y, Z float32
)

func main() {
	fmt.Printf("X size:%v, Y size:%v, Z size:%v\n", unsafe.Sizeof(X), unsafe.Sizeof(Y), unsafe.Sizeof(Z))
	fmt.Printf("X addr:%v, Y addr:%v, Z addr:%v\n", &X, &Y, &Z)
}
