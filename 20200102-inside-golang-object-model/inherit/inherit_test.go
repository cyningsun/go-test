package inherit

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

func TestLayout(t *testing.T) {
	p := point3d{point2d: point2d{point: point{X: 1}, Y: 2}, Z: 3}
	fmt.Printf("point3d size:%v, align:%v\n", unsafe.Sizeof(p), unsafe.Alignof(p))
	typ := reflect.TypeOf(p)
	fmt.Printf("Struct:%v is %d bytes long\n", typ.Name(), typ.Size())
	fmt.Printf("X at offset %v, size=%d\n", unsafe.Offsetof(p.X), unsafe.Sizeof(p.X))
	fmt.Printf("Y at offset %v, size=%d\n", unsafe.Offsetof(p.Y), unsafe.Sizeof(p.Y))
	fmt.Printf("Z at offset %v, size=%d\n", unsafe.Offsetof(p.Z), unsafe.Sizeof(p.Z))
}

func TestPolymorphism(t *testing.T) {
	var (
		p    Point
		nilP interface{}
	)
	p = &point{X: 1}
	nilP = &point{X: 1}
	fmt.Printf("p size:%v, nilP size:%v\n", unsafe.Sizeof(p), unsafe.Sizeof(nilP))

	p = &point2d{point: point{X: 1}, Y: 2}
	nilP = &point2d{point: point{X: 1}, Y: 2}
	fmt.Printf("p size:%v, nilP size:%v\n", unsafe.Sizeof(p), unsafe.Sizeof(nilP))

	p = &point3d{point2d: point2d{point: point{X: 1}, Y: 2}, Z: 3}
	nilP = &point3d{point2d: point2d{point: point{X: 1}, Y: 2}, Z: 3}
	fmt.Printf("p size:%v, nilP size:%v\n", unsafe.Sizeof(p), unsafe.Sizeof(nilP))
}
