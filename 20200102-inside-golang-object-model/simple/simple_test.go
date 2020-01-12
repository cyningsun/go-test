package simple

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"

	"github.com/cyningsun/go-test/20200102-inside-golang-object-model/face"
	"github.com/davecgh/go-spew/spew"
)

func TestLayout(t *testing.T) {
	p := point3d{X: 1, Y: 2, Z: 3}
	fmt.Printf("point3d size:%v, align:%v\n", unsafe.Sizeof(p), unsafe.Alignof(p))
	typ := reflect.TypeOf(p)
	fmt.Printf("Struct:%v is %d bytes long\n", typ.Name(), typ.Size())
	fmt.Printf("X at offset %v, size=%d\n", unsafe.Offsetof(p.X), unsafe.Sizeof(p.X))
	fmt.Printf("Y at offset %v, size=%d\n", unsafe.Offsetof(p.Y), unsafe.Sizeof(p.Y))
	fmt.Printf("Z at offset %v, size=%d\n", unsafe.Offsetof(p.Z), unsafe.Sizeof(p.Z))
}

func TestInterface(t *testing.T) {
	var (
		p    Point
		nilP interface{}
	)
	point := &point3d{X: 1, Y: 2, Z: 3}
	nilP = point
	fmt.Printf("eface size:%v\n", unsafe.Sizeof(nilP))
	eface := (*face.Eface)(unsafe.Pointer(&nilP))
	spew.Dump(eface.Type)
	spew.Dump(eface.Data)
	fmt.Printf("eface offset: eface._type = %v, eface.data = %v\n\n",
		unsafe.Offsetof(eface.Type), unsafe.Offsetof(eface.Data))

	p = point
	fmt.Printf("point size:%v\n", unsafe.Sizeof(p))
	iface := (*face.Iface)(unsafe.Pointer(&p))
	spew.Dump(iface.Tab)
	spew.Dump(iface.Data)
	fmt.Printf("Iface offset: iface.tab = %v, iface.data = %v\n\n",
		unsafe.Offsetof(iface.Tab), unsafe.Offsetof(iface.Data))
}
