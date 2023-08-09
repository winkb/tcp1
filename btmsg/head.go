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
	defer func() {
		if err != nil {
			fmt.Println(err)
		}
	}()

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

	err = binary.Read(bytes.NewReader(hdBt), binary.LittleEndian, l)
	if err != nil {
		err = errors.Wrap(err, "byte to head")
		return
	}

	return nil
}

func (l *MsgHead) ToBytes() []byte {
	bt := make([]byte, 0)
	bf := bytes.NewBuffer(bt)
	_ = binary.Write(bf, binary.LittleEndian, l)
	return bf.Bytes()
}
