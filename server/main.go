package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"tcp1/btmsg"
	"tcp1/mytcp"
	"tcp1/util/numfn"
	"time"
)

type RouteHandle func(conn *mytcp.TcpConn, msg btmsg.IMsg, req any)

type RouteInfo struct {
	Handle RouteHandle
	Info   any
}

var server mytcp.ITcpServer

type ShutdownReq struct {
	Msg string
}

type ShutdownRsp struct {
	Reason string
}

var routes = map[uint16]*RouteInfo{}

func init() {
	routes[100] = &RouteInfo{
		Handle: func(conn *mytcp.TcpConn, msg btmsg.IMsg, req any) {
			handleShutdown(conn, msg, req.(*ShutdownReq))
		},
		Info: &ShutdownReq{},
	}
}

func logHandle(name string, t time.Time) func() {
	return func() {
		fmt.Println("handle", name, "in")
		fmt.Println("handle", name, "out", numfn.ToStr(time.Now().Sub(t).Nanoseconds())+"ns")
	}
}

func handleShutdown(conn *mytcp.TcpConn, msg btmsg.IMsg, req *ShutdownReq) {
	defer logHandle("shutdown", time.Now())

	fmt.Println("sever will shutdown ", req.Msg)

	err := msg.FromStruct(&ShutdownRsp{
		Reason: "server will shutdown! trigger by " + conn.GetRemoteIp(),
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	server.Broadcast(msg)

	server.Shutdown()
}

func main() {
	server = mytcp.NewTcpServer("989", btmsg.NewReader())
	wg, err := server.Start()
	if err != nil {
		panic(err)
	}

	server.OnReceive(func(conn *mytcp.TcpConn, msg btmsg.IMsg) {
		act := msg.GetAct()
		hv, ok := routes[act]
		if !ok {
			fmt.Println("not found handle", act)
			return
		}

		var info = hv.Info
		info, err = msg.ToStruct(hv.Info)
		if err != nil {
			fmt.Println(err)
			return
		}

		hv.Handle(conn, msg, info)
	})

	chSingle := make(chan os.Signal)

	signal.Notify(chSingle, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		v := <-chSingle
		switch v {
		case syscall.SIGINT:
			fmt.Println("ctr+c")
		case syscall.SIGTERM:
			fmt.Println("terminated")
		}

		server.Shutdown()

		fmt.Println(v)
	}()

	wg.Wait()
	close(chSingle)
}
