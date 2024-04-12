package mytcp

import (
	"net"
)

type wrapConn struct {
	net.Conn
}

func NewWrapConn(conn net.Conn)*wrapConn  {
	return &wrapConn{conn}
}

func (l *wrapConn) ReadMessage() (messageType int, p []byte, err error) {
	//TODO implement me
	panic("implement me")
}

func (l *wrapConn) GetRemoteIp() string {
	return l.Conn.RemoteAddr().String()
}

