package main

import (
	"fmt"
	"sync"
	"tcp1/mytcp"
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
			v.Send("shutdown;")
		}
	})

	time.Sleep(time.Second * 10)
}

func start() (mytcp.ITcpClient, *sync.WaitGroup) {
	cli := mytcp.NewTcpClient(":989")

	cli.OnReceive(func(v []byte) {
		if string(v) == "panic;" {
			panic("panic by user")
			return
		}
		fmt.Println("服务端回复:", string(v))
	})
	cli.OnClose(func(isServer bool, isClient bool) {
		if isClient {
			fmt.Println("客户端断开连接")
		}

		if isServer {
			cli.ReleaseChan()
			fmt.Println("服务端断开连接")
		}
	})

	wg, err := cli.Start()
	if err != nil {
		panic(err)
	}

	return cli, wg
}
