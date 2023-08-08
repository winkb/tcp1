package btmsg

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"unsafe"
)

type MsgHead struct {
	Act  uint16
	Size uint32
}

var _ IHead = (*MsgHead)(nil)

func NewMsgHead() *MsgHead {
	return &MsgHead{}
}

func (l *MsgHead) GetAct() uint16 {
	return l.Act
}

func (l *MsgHead) HeadSize() uint32 {
	return uint32(unsafe.Sizeof(l.Act) + unsafe.Sizeof(l.Size))
}

func (l *MsgHead) BodySize() uint32 {
	return l.Size
}

func (l *MsgHead) Read(r io.Reader) (err error) {
	var headSize = l.HeadSize()
	var n int
	var hdBt = make([]byte, headSize)

	n, err = r.Read(hdBt)
	if err != nil {
		return err
	}

	if n != int(headSize) {
		return fmt.Errorf("head len err got %v, expect %v", n, headSize)
	}

	err = binary.Read(bytes.NewReader(hdBt), binary.LittleEndian, l)
	if err != nil {
		return errors.Wrap(err, "byte to head")
	}

	return nil
}

func (l *MsgHead) ReadBody(r io.Reader) (err error) {
	var headSize = l.BodySize()
	var n int
	var hdBt = make([]byte, headSize)

	n, err = r.Read(hdBt)
	if err != nil {
		return err
	}

	if n != int(headSize) {
		return fmt.Errorf("head len err got %v, expect %v", n, headSize)
	}

	err = binary.Read(bytes.NewReader(hdBt), binary.LittleEndian, l)
	if err != nil {
		return errors.Wrap(err, "byte to head")
	}

	return nil
}

func (l *MsgHead) ToBytes() (bt []byte) {
	var act = byte(l.Act)
	var size = byte(l.Size)
	return append(bt, act, size)
}