package main

import "fmt"

type point3d struct {
	X, Y, Z float32
}

func (p *point3d) Println() {
	fmt.Printf("%v,%v,%v\n", p.X, p.Y, p.Z)
}

func main() {
	p := point3d{X: 1, Y: 2, Z: 3}
	p.Println()
}
