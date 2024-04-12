package handles

import (
	"github.com/winkb/tcp1/btmsg"
	"github.com/winkb/tcp1/cmd/server/types"
	"github.com/winkb/tcp1/contracts"
)

func parseReq[T any](v T, msg btmsg.IMsg) T {
	x,_ := msg.ToStruct(v)
	return x.(T)
}

type RouteHandle func(s contracts.ITcpServer,conn *contracts.TcpConn, msg btmsg.IMsg)

type RouteInfo struct {
	Handle RouteHandle
}
var Routes = map[uint16]*RouteInfo{}

func init() {
	Routes[0] = &RouteInfo{
		Handle: func(s contracts.ITcpServer,conn *contracts.TcpConn, msg btmsg.IMsg) {
			handleDefault(s,conn, msg, nil)
		},
	}

	Routes[1] = &RouteInfo{
		Handle: func(s contracts.ITcpServer,conn *contracts.TcpConn, msg btmsg.IMsg) {
			handleHello(s,conn, msg, parseReq(&types.HelloReq{}, msg))
		},
	}

	Routes[100] = &RouteInfo{
		Handle: func(s contracts.ITcpServer,conn *contracts.TcpConn, msg btmsg.IMsg) {
			handleShutdown(s,conn, msg, parseReq(&types.ShutdownReq{}, msg))
		},
	}
}
