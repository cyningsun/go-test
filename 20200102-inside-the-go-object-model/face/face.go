package face

import "unsafe"

type Iface struct {
	Tab *Itab

	Data unsafe.Pointer
}

type Itab struct {
	Inter uintptr

	Type uintptr

	Hash uint32

	_ [4]byte

	Fun [1]uintptr
}

type Eface struct {
	Type uintptr

	Data unsafe.Pointer
}
