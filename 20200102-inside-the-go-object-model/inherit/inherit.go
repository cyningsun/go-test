package inherit

import "fmt"

type Point interface {
	Println()
}

type point struct {
	X float32
}

func (p *point) Println() {
	fmt.Printf("%v\n", p.X)
}

type point2d struct {
	point
	Y float32
}

func (p *point2d) Println() {
	fmt.Printf("%v,%v\n", p.X, p.Y)
}

type point3d struct {
	point2d
	Z float32
}

func (p *point3d) Println() {
	fmt.Printf("%v,%v,%v\n", p.X, p.Y, p.Z)
}
