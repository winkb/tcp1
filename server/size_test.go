package main

import (
	"testing"
	"unsafe"
)

type head1 struct {
	Act  uint8
	Size uint16
}

func TestSize1(t *testing.T) {
	var i1 uint16 = 1

	t.Log(unsafe.Sizeof(i1))
}

func TestSize2(t *testing.T) {
	var h1 = &head1{}

	t.Log(unsafe.Sizeof(h1.Act))
	t.Log(unsafe.Sizeof(h1.Size))
}
