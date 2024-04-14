package contracts

import (
	"net"
	"sync"

	"github.com/winkb/tcp1/btmsg"
)

type ServerCloseCallback func(s ITcpServer, conn *TcpConn, isServer bool, isClient bool)
type ServerReceiveCallback func(s ITcpServer, conn *TcpConn, msg btmsg.IMsg)

type ITcpServer interface {
	Shutdown()
	Send(conn *TcpConn, v btmsg.IMsg)
	Close(conn *TcpConn)
	SendById(id uint64, v btmsg.IMsg)
	OnReceive(f ServerReceiveCallback)
	OnClose(f ServerCloseCallback)
	Start() (wg *sync.WaitGroup, err error)
	Broadcast(bt btmsg.IMsg)
}

type IConn interface {
	GetRemoteIp() string
	net.Conn
	ReadMessage() (messageType int, p []byte, err error)
}

type TcpConn struct {
	Conn     IConn
	Id       uint64
	Input    chan btmsg.IMsg
	Output   chan btmsg.IMsg
	WaitConn chan bool
	Lock     sync.RWMutex
	IsClose  bool
}

func (l *TcpConn) GetRemoteIp() string {
	if l.Conn == nil {
		return ""
	}
	return l.Conn.GetRemoteIp()
}

func (l *TcpConn) GetId() uint64 {
	return l.Id
}
