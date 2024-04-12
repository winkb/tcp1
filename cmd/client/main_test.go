package main

import (
	"fmt"
	"github.com/winkb/tcp1/btmsg"
	"github.com/winkb/tcp1/net/mytcp"
	"sync"
	"testing"
	"time"
)

func TestShutdown(t *testing.T) {
	var conns []mytcp.ITcpClient
	var lock sync.Mutex

	go func() {
		for i := 0; i < 100; i++ {
			cli, _ := start()

			lock.Lock()
			conns = append(conns, cli)
			lock.Unlock()
		}
	}()

	time.AfterFunc(time.Second*3, func() {
		for _, v := range conns {
			msg := newMsg(100, &ShutdownReq{
				Msg: "exit",
			})
			v.Send(msg)
		}
	})

	time.Sleep(time.Second * 10)
}

func start() (mytcp.ITcpClient, *sync.WaitGroup) {
	cli := mytcp.NewTcpClient(":989")

	cli.OnReceive(func(v btmsg.IMsg) {
		fmt.Println(v.GetAct(), string(v.BodyByte()))
	})
	cli.OnClose(func(isServer bool, isClient bool) {
		if isClient {
			fmt.Println("服务端断开连接")
		}

		if isServer {
			cli.ReleaseChan()
			fmt.Println("我自己断开连接")
		}
	})

	wg, err := cli.Start()
	if err != nil {
		panic(err)
	}

	return cli, wg
}
