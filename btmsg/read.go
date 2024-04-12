package btmsg

type Reader struct {
	f func()IHead
}

func NewReader(f func()IHead) *Reader {
	return &Reader{
		f: f,
	}
}

func (l *Reader) ReadMsg(r IReader) (res IReadResult) {
	var err error
	var head = l.f()

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
