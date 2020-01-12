package simple

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
