package btmsg

import (
	"io"
)

type IHead interface {
	BodySize() uint32
	HeadSize() uint32
	GetAct() uint16
}

type IMsg interface {
	IHead
	BodyByte() []byte
	FromStruct(v any) error
	ToStruct(v any) (any, error)
	ToSendByte() []byte
}

type IReadResult interface {
	IsClose() bool
	IsCloseByServer() bool
	IsCloseByClient() bool
	GetMsg() IMsg
	GetErr() error
}

type IMsgReader interface {
	ReadMsg(r io.Reader) (res IReadResult)
}
