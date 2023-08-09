package btmsg

import (
	"fmt"
	"io"
)

type Reader struct {
}

func NewReader() *Reader {
	return &Reader{}
}

func (l *Reader) ReadMsg(r io.Reader) (res IReadResult) {
	var err error
	var head = NewMsgHead()

	{
		err = head.Read(r)
		if err != nil {
			return NewReaerResult(err, head, nil)
		}
	}

	var bodySize = head.BodySize()
	var bt = make([]byte, bodySize)

	{
		var n int
		n, err = io.ReadFull(r, bt)
		if err != nil {
			return NewReaerResult(err, head, nil)
		}

		if n != int(bodySize) {
			return NewReaerResult(fmt.Errorf("body len err got %v, expect %v", n, bodySize), head, nil)
		}
	}

	return NewReaerResult(err, head, bt)
}
