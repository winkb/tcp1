package btmsg

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"unsafe"
)

// MsgHeadTcp 作为容器，不可以有多余字段
type MsgHeadTcp struct {
	Act  uint16
	Size uint32
}

var _ IHead = (*MsgHeadTcp)(nil)

func FactoryMsgHeadTcp() func() IHead {
	return func() IHead {
		return NewMsgHeadTcp()
	}
}

func NewMsgHeadTcp() *MsgHeadTcp {
	return &MsgHeadTcp{}
}

func (l *MsgHeadTcp) SetSize(size uint32) {
	l.Size = size
}

func (l *MsgHeadTcp) GetAct() uint16 {
	return l.Act
}

// 如果MsgHead增加了字段，这里也要对应修改
func (l *MsgHeadTcp) HeadSize() uint32 {
	return uint32(unsafe.Sizeof(l.Act) + unsafe.Sizeof(l.Size))
}

func (l *MsgHeadTcp) BodySize() uint32 {
	return l.Size
}

func (l *MsgHeadTcp) Read(r IReader) (err error) {
	var headSize = l.HeadSize()
	var n int
	var hdBt = make([]byte, headSize)

	n, err = io.ReadFull(r, hdBt)
	if err != nil {
		return err
	}

	if n != int(headSize) {
		err = fmt.Errorf("head len err got %v, expect %v", n, headSize)
		return
	}

	// 将head 字节 解析成结构体
	err = binary.Read(bytes.NewReader(hdBt), binary.LittleEndian, l)
	if err != nil {
		err = errors.Wrap(err, "byte to head")
		return
	}

	return nil
}

func (l *MsgHeadTcp) ReadBody(r IReader) (err error, bt []byte) {
	var bodySize = l.BodySize()
	var n int

	bt = make([]byte, bodySize)
	n, err = io.ReadFull(r, bt)
	if err != nil {
		return
	}

	if n != int(bodySize) {
		err = fmt.Errorf("body len err got %v, expect %v", n, bodySize)
		return
	}

	return
}

func (l *MsgHeadTcp) ToBytes() []byte {
	bt := make([]byte, 0)
	bf := bytes.NewBuffer(bt)
	_ = binary.Write(bf, binary.LittleEndian, l)
	return bf.Bytes()
}

func (l *MsgHeadTcp) FromStruct(v any) (bt []byte, err error) {
	bt, err = json.Marshal(&WsResponse[any]{
		Act: l.GetAct(),
		Data: v,
	})
	if err != nil {
		err = errors.Wrap(err, "struct to msg")
		return
	}

	return
}

func (l *MsgHeadTcp) ToStruct(bt []byte, v any) (any, error) {
	err := json.Unmarshal(bt, v)
	if err != nil {
		return v, errors.Wrap(err, "msg to struct")
	}

	return v, nil
}


func (l *MsgHeadTcp) SetAct(act uint16)  {
	l.Act = act
}
