package main

import "fmt"

type Point interface {
	Println()
}

type point3d struct {
	X, Y, Z float32
}

func (p *point3d) Println() {
	fmt.Printf("%v,%v,%v\n", p.X, p.Y, p.Z)
}

func main() {
	point := point3d{X: 1, Y: 2, Z: 3}
	var (
		nilP interface{}
		p    Point
	)
	nilP = &point
	p = &point
	fmt.Println(nilP, p)
}

/*
 * go build -gcflags '-N -l' -o build build.go
`* go tool objdump -s "main\.main" build
*/
