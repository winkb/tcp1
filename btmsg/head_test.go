package btmsg

import (
	"bytes"
	"encoding/binary"
	"testing"
	"unsafe"
)

func TestBtToHead(t *testing.T) {
	var h = &MsgHead{}
	var bt = make([]byte, h.HeadSize())

	t.Log(unsafe.Sizeof(h.Act))
	t.Log(unsafe.Sizeof(h.Size))

	bt[0] = 100
	bt[2] = 100

	t.Log(bt)

	err := binary.Read(bytes.NewReader(bt), binary.LittleEndian, h)
	if err != nil {
		panic(err)
	}

	t.Logf("%v", h)
}
