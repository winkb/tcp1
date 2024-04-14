package btmsg

var _ IMsg = (*Msg)(nil)

type Msg struct {
	head   IHead
	bodyBt []byte
}

func (l *Msg) BodySize() uint32 {
	return l.head.BodySize()
}

func (l *Msg) SetAct(act uint16) {
	l.head.SetAct(act)
}
func (l *Msg) HeadSize() uint32 {
	return l.head.HeadSize()
}

func (l *Msg) GetAct() uint16 {
	return l.head.GetAct()
}

func NewMsg(head IHead, bodyBt []byte) *Msg {
	return &Msg{
		head:   head,
		bodyBt: bodyBt,
	}
}

func (l *Msg) BodyByte() []byte {
	return l.bodyBt
}

// v is a pointer
func (l *Msg) FromStruct(v any) (err error) {
	l.bodyBt, err = l.head.FromStruct(v)
	return
}

// v is a pointer
func (l *Msg) ToStruct(v any) (any, error) {
	return l.head.ToStruct(l.bodyBt, v)
}

// 除非只需要发送head,否则需要在FromStruct之后执行
func (l *Msg) ToSendByte() []byte {
	l.head.SetSize(uint32(len(l.bodyBt)))

	bt := l.head.ToBytes()
	bt = append(bt, l.bodyBt...)

	return bt
}
