package btmsg

import (
	"bytes"
	"encoding/binary"
	"github.com/pkg/errors"
)

var _ IMsg = (*Msg)(nil)

type Msg struct {
	*MsgHead
	bodyBt []byte
}

func NewMsg(head *MsgHead, bodyBt []byte) *Msg {
	return &Msg{
		MsgHead: head,
		bodyBt:  bodyBt,
	}
}

func (l *Msg) BodyByte() []byte {
	return l.bodyBt
}

func (l *Msg) FromStruct(v any) (err error) {
	var bt = make([]byte, l.BodySize())
	var w = bytes.NewBuffer(bt)
	err = binary.Write(w, binary.LittleEndian, v)
	if err != nil {
		err = errors.Wrap(err, "struct to msg")
		return
	}

	return nil
}

func (l *Msg) ToStruct(v any) (any, error) {
	var bt = make([]byte, l.BodySize())
	err := binary.Read(bytes.NewReader(bt), binary.LittleEndian, v)
	if err != nil {
		return v, errors.Wrap(err, "msg to struct")
	}

	return v, nil
}

func (l *Msg) ToByte() []byte {
	bt := l.MsgHead.ToBytes()
	bt = append(bt, l.bodyBt...)
	return bt
}
