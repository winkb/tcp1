package myws

import (
	"time"

	"github.com/gorilla/websocket"
)

type wrapConn struct {
	*websocket.Conn
}

func (l *wrapConn) ReadMessage() (messageType int, p []byte, err error) {
	messageType, p, err = l.Conn.ReadMessage()
	return
}

func (l *wrapConn) Read(b []byte) (n int, err error) {
	_, b, err = l.Conn.ReadMessage()
	n  = len(b)
	return
}

func (l *wrapConn) SetDeadline(t time.Time) error {
	if err := l.Conn.SetReadDeadline(t); err != nil {
		return err
	}
	return l.Conn.SetWriteDeadline(t)
}

func (l *wrapConn) GetRemoteIp() string {
	return l.Conn.RemoteAddr().String()
}

// Write writes data to the connection.
// Write can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetWriteDeadline.
func (l *wrapConn) Write(b []byte) (n int, err error) {
	return
}

