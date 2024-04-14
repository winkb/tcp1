package btmsg

type IHead interface {
	BodySize() uint32
	HeadSize() uint32
	GetAct() uint16
	SetSize(size uint32)
	ToBytes() []byte
	Read(r IReader) (err error)
	ReadBody(r IReader) (err error, bt []byte)
	FromStruct(v any) (bt []byte, err error)
	ToStruct(bt []byte, v any) (any, error)
	SetAct(act uint16)
}

type IReader interface {
	ReadMessage() (messageType int, p []byte, err error)
	Read(b []byte) (n int, err error)
}

type IMsg interface {
	BodySize() uint32
	HeadSize() uint32
	GetAct() uint16
	BodyByte() []byte
	FromStruct(v any) error
	ToStruct(v any) (any, error)
	ToSendByte() []byte
	SetAct(act uint16)
}

type IReadResult interface {
	IsClose() bool
	IsCloseByServer() bool
	IsCloseByClient() bool
	GetMsg() IMsg
	GetErr() error
}

type IMsgReader interface {
	ReadMsg(r IReader) (res IReadResult)
}
