package main

import (
	"bufio"
	"fmt"
	"github.com/winkb/tcp1/btmsg"
	"github.com/winkb/tcp1/mytcp"
	"os"
)

type RouteHandle func(msg btmsg.IMsg)

type RouteInfo struct {
	Handle RouteHandle
}

type ShutdownReq struct {
	Msg string
}

type ShutdownRsp struct {
	Reason string
}

func newMsg(act uint16, req any) btmsg.IMsg {
	hd := btmsg.NewMsgHead()
	hd.Act = act
	res := btmsg.NewMsg(hd, nil)
	err := res.FromStruct(req)
	if err != nil {
		fmt.Println(err)
	}
	return res
}

var routes = map[uint16]*RouteInfo{}

func init() {
	routes[100] = &RouteInfo{
		Handle: func(msg btmsg.IMsg) {
			handleShutdownReply(msg, parseReq(ShutdownRsp{}, msg))
		},
	}
}

func parseReq[T any](v T, msg btmsg.IMsg) T {
	_, _ = msg.ToStruct(&v)
	return v
}

func handleShutdownReply(msg btmsg.IMsg, req ShutdownRsp) {
	fmt.Println("shutdown notify ", req.Reason)
}

func main() {

	cli := mytcp.NewTcpClient(":989")

	cli.OnReceive(func(v btmsg.IMsg) {
		act := v.GetAct()
		r, ok := routes[act]
		if !ok {
			fmt.Println("not found handle", act)
			return
		}

		r.Handle(v)
	})

	cli.OnClose(func(isServer bool, isClient bool) {
		if isClient {
			fmt.Println("服务端断开连接")
		}

		if isServer {
			cli.ReleaseChan()
			fmt.Println("我自己端口连接")
		}
	})

	wg, err := cli.Start()
	if err != nil {
		panic(err)
	}

	scan := bufio.NewScanner(os.Stdin)
	const exitLimit = "exit;"

	mytcp.MyGoWg(wg, "scan_input", func() {
		defer func() {
			cli.Close()
		}()
		for scan.Scan() {

			txt := scan.Text()
			if txt == exitLimit {
				return
			}

			select {
			case <-cli.HasClosed():
				return
			default:
				if txt == "shutdown" {
					cli.Send(newMsg(100, ShutdownReq{
						Msg: txt,
					}))
					continue
				}

				cli.Send(newMsg(200, ShutdownReq{
					Msg: txt,
				}))
			}
		}
	})

	wg.Wait()
}
