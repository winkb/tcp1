package btmsg

import (
	"io"
)

type IHead interface {
	BodySize() uint32
}

type IMsg interface {
	IHead
	Byte() []byte
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
