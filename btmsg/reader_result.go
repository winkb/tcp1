package btmsg

import (
	"io"
	"net"
)

type ReaderResult struct {
	err  error
	head *MsgHead
	body []byte
}

func NewReaerResult(err error, head *MsgHead, body []byte) *ReaderResult {
	return &ReaderResult{
		err:  err,
		head: head,
		body: body,
	}
}

func (l *ReaderResult) IsClose() bool {
	return l.err != nil
}

func (l *ReaderResult) IsCloseByServer() bool {
	if l.err == nil {
		return false
	}
	_, ok := l.err.(*net.OpError)
	return ok
}

func (l *ReaderResult) IsCloseByClient() bool {
	return l.err == io.EOF
}

func (l *ReaderResult) GetErr() error {
	return l.err
}

func (l *ReaderResult) GetMsg() IMsg {
	return NewMsg(l.head, l.body)
}
