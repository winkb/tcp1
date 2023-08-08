package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"tcp1/mytcp"
)

func main() {
	server := mytcp.NewTcpServer("989")
	wg, err := server.Start()
	if err != nil {
		panic(err)
	}

	server.OnReceive(func(conn *mytcp.TcpConn, bt []byte) {
		if string(bt) == "shutdown;" {
			fmt.Println("receive: shutdown")
			server.Shutdown()
		}

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
