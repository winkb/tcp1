package btmsg

import (
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

	err = head.Read(r)
	if err != nil {
		return NewReaderResult(err, head, nil)
	}

	var body []byte
	err, body = head.ReadBody(r)
	if err != nil {
		return NewReaderResult(err, head, nil)
	}

	return NewReaderResult(err, head, body)
}

type ReaderWs struct {
}

func NewReaderWs() *ReaderWs {
	return &ReaderWs{}
}

func (l *ReaderWs) ReadMsg(r io.Reader) (res IReadResult) {
	var err error
	var head = NewMsgHead()

	err = head.Read(r)
	if err != nil {
		return NewReaderResult(err, head, nil)
	}

	var body []byte
	err, body = head.ReadBody(r)
	if err != nil {
		return NewReaderResult(err, head, nil)
	}

	return NewReaderResult(err, head, body)
}
