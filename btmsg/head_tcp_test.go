package btmsg

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"testing"
	"unsafe"
)

func TestBtToHead(t *testing.T) {
	var h = &MsgHeadTcp{}
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

func TestHeadBytes2(t *testing.T) {
	var h = &MsgHeadTcp{}

	h.Act = 100
	h.Size = 12

	bt := h.ToBytes()
	t.Log(bt)
}

func TestWriteHead(t *testing.T) {
	var h = &MsgHeadTcp{
		Act:  100,
		Size: 12,
	}
	var bt = make([]byte, 0)
	var bf = bytes.NewBuffer(bt)
	err := binary.Write(bf, binary.LittleEndian, h)
	t.Log(err)

	t.Log(bf.Bytes())
}

type req1 struct {
	Msg uint32
}

func TestWriteStructToBinary(t *testing.T) {
	var r = req1{
		Msg: 10,
	}

	marshal, err2 := json.Marshal(r)
	t.Log(err2)
	t.Log(marshal)

	// 下面是错误的方法
	var bt = make([]byte, 0)
	var bf = bytes.NewBuffer(bt)
	err := binary.Write(bf, binary.LittleEndian, r)
	t.Log(err)
	t.Log(bf.Bytes())
}
