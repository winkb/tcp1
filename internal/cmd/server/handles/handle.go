package handles

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/winkb/tcp1/btmsg"
	"github.com/winkb/tcp1/internal/cmd/server/types"
	"github.com/winkb/tcp1/contracts"
	"github.com/winkb/tcp1/util/numfn"
	"time"
)


func handleDefault(s contracts.ITcpServer,conn *contracts.TcpConn, msg btmsg.IMsg, req any) {
	defer logHandle("default", time.Now())

	fmt.Println("sever receive default msg ", req)
}

func handleShutdown(s contracts.ITcpServer,conn *contracts.TcpConn, msg btmsg.IMsg, req *types.ShutdownReq) {
	defer logHandle("shutdown", time.Now())

	fmt.Println("sever will shutdown ", req.Msg)

	err := msg.FromStruct(&types.ShutdownRsp{
		Reason: "server will shutdown! trigger by " + conn.GetRemoteIp(),
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	s.Broadcast(msg)
	time.AfterFunc(time.Second, func() {
		s.Shutdown()
	})
}

func handleHello(s contracts.ITcpServer,conn *contracts.TcpConn, msg btmsg.IMsg, req *types.HelloReq) {
	fmt.Println("hello", req.Content)

	err := msg.FromStruct(req)
	if err != nil {
		log.Err(err)
		return
	}

	s.Send(conn, msg)
}

func logHandle(name string, t time.Time) func() {
	return func() {
		fmt.Println("handle", name, "in")
		fmt.Println("handle", name, "out", numfn.ToStr(time.Now().Sub(t).Nanoseconds())+"ns")
	}
}
