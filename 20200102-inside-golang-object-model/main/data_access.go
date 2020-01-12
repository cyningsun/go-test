package main

import (
	"fmt"
)

type point struct {
	X float32
}

type point2d struct {
	point
	Y float32
}

type point3d struct {
	point2d
	Z float32
}

func main() {
	var (
		w float32
	)
	point := point3d{point2d: point2d{point: point{X: 1}, Y: 2}, Z: 3}
	p := &point
	w = point.Y
	fmt.Printf("w:%f\n", w)
	w = p.Y
	fmt.Printf("w:%f\n", w)
}

/*
 * go build -gcflags '-N -l' -o data_access data_access.go
`* go tool objdump -s "main\.main" data_access
*/
