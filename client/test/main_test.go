package main

import (
	"bufio"
	"fmt"
	"github.com/winkb/tcp1/mytcp"
	"os"
	"sync"
	"testing"
	"time"
)

func startClient() {
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
				cli.Send(txt)
			}
		}
	})

	wg.Wait()
}

func TestClientMul(t *testing.T) {
	wg := sync.WaitGroup{}

	const N = 100

	wg.Add(N)
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()

			startClient()
		}()

	}

	wg.Wait()

	time.Sleep(time.Second * 6)
}
