package main

import (
	"fmt"
	"github.com/winkb/tcp1/contracts"
	"github.com/winkb/tcp1/util"
	"sync"
	"testing"
	"time"

	"github.com/winkb/tcp1/btmsg"
	"github.com/winkb/tcp1/net/mytcp"
)

func startClient() {

	cli := mytcp.NewTcpClient(":989")

	cli.OnReceive(func(msg btmsg.IMsg) {
		var v = msg.BodyByte()

		if string(v) == "panic;" {
			panic("panic by user")
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

	const exitLimit = "exit;"

	util.MyGoWg(wg, "scan_input", func() {
		defer func() {
			cli.Close()
		}()
	})

	wg.Wait()
}

func TestClientMul(t *testing.T) {
	wg := sync.WaitGroup{}

	const N = 100

	ts := mytcp.NewTcpServer("989", btmsg.NewReader(func() btmsg.IHead {
		return btmsg.NewMsgHeadTcp()
	}))

	var num int32 = 0
	var lock sync.RWMutex

	go func() {

		wg2, err2 := ts.Start()
		if err2 != nil {
			panic(err2)
		}

		ts.OnClose(func(s contracts.ITcpServer,conn *contracts.TcpConn, isServer bool, isClient bool) {
			if isClient {
				lock.Lock()
				num++

				if num >= N {
					ts.Shutdown()
				}

				lock.Unlock()

				t.Log("conn close by client self")
			}
		})

		wg2.Wait()

		t.Log("close_all_conn", num)
	}()

	time.Sleep(time.Second)

	wg.Add(N)
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()

			startClient()
		}()

	}

	wg.Wait()

	time.Sleep(time.Second * 2)
}
