package btmsg

import (
	"bytes"
	"io"
	"testing"
)

type mockReader struct{
	reader *bytes.Reader
}

func (l *mockReader) ReadMessage() (messageType int, p []byte, err error) {
	p, err = io.ReadAll(l.reader)
	return  messageType, p, err
}

func (l *mockReader) Read(b []byte) (n int, err error) {
	//TODO implement me
	panic("implement me")
}

func TestMsgHeadWs_Read(t *testing.T) {
	var rd = bytes.NewReader([]byte(`{"act":1,"content":"abc"}`))

	ws := NewMsgHeadWs()

	err := ws.Read(&mockReader{rd})
	t.Log(err)

	act := ws.GetAct()
	t.Log(act)
}
