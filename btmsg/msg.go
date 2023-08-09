package btmsg

import (
	"encoding/json"
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

// v is a pointer
func (l *Msg) FromStruct(v any) (err error) {
	var bt []byte
	bt, err = json.Marshal(v)
	if err != nil {
		err = errors.Wrap(err, "struct to msg")
		return
	}

	l.bodyBt = bt

	return nil
}

// v is a pointer
func (l *Msg) ToStruct(v any) (any, error) {
	err := json.Unmarshal(l.bodyBt, v)
	if err != nil {
		return v, errors.Wrap(err, "msg to struct")
	}

	return v, nil
}

// 除非只需要发送head,否则需要在FromStruct之后执行
func (l *Msg) ToSendByte() []byte {
	l.MsgHead.Size = uint32(len(l.bodyBt))

	bt := l.MsgHead.ToBytes()
	bt = append(bt, l.bodyBt...)

	return bt
}
