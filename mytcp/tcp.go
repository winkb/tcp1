package mytcp

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"net"
	"sync"
)

func MyGoWg(wg *sync.WaitGroup, name string, f func()) {
	wg.Add(1)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Print(fmt.Sprintf("goroutine %s %s", name, err))

				MyGoWg(wg, name, f)

				wg.Done()
				return
			}

			fmt.Printf("goroutine %s defer\n", name)

			wg.Done()
		}()

		f()
	}()
}

func myGo(name string, f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Print(name, err)
				return
			}

			fmt.Printf("goroutine %s defer", name)
		}()

		f()
	}()
}

type ITcpClient interface {
	LoopRead()
	ReleaseChan()
	LoopWrite()
	LoopReceive()
	Close()
	Send(v string)
	OnReceive(f func(v []byte))
	OnClose(f func(isServer bool, isClient bool))
	Start() (wg *sync.WaitGroup, err error)
	HasClosed() chan bool
}

var _ ITcpClient = (*tcpClient)(nil)

type tcpClient struct {
	input           chan string
	output          chan []byte
	wait            chan bool
	conn            net.Conn
	closeCallback   func(isServer bool, isClient bool)
	receiveCallback func(bt []byte)
	addr            string
}

func (l *tcpClient) Start() (wg *sync.WaitGroup, err error) {
	wg = &sync.WaitGroup{}
	// conn server
	err = l.connServer()
	if err != nil {
		return
	}
	// read
	MyGoWg(wg, "conn_read", l.LoopRead)
	// write
	MyGoWg(wg, "conn_write", l.LoopWrite)
	// on msg
	MyGoWg(wg, "conn_receive", l.LoopReceive)

	return
}

func (l *tcpClient) connServer() error {
	conn, err := net.Dial("tcp", l.addr)
	if err != nil {
		return err
	}
	l.conn = conn
	return nil
}

func (l *tcpClient) handelReadClose(isServer bool, isClient bool) {
	close(l.wait)
	if l.closeCallback != nil {
		l.closeCallback(isServer, isClient)
	}
}

func (l *tcpClient) handelReceive(bt []byte) {
	if l.receiveCallback != nil {
		l.receiveCallback(bt)
	}
}

func (l *tcpClient) OnClose(f func(isServer bool, isClient bool)) {
	l.closeCallback = f
}

func (l *tcpClient) LoopRead() {
	for {
		bt := make([]byte, 1024)
		n, err := l.conn.Read(bt)
		if err == io.EOF {
			l.handelReadClose(true, false)
			return
		}

		if _, ok := err.(*net.OpError); ok {
			l.handelReadClose(false, true)
			return
		}

		if err != nil {
			l.log("conn read", err)
			continue
		}

		msg := bt[:n]
		l.output <- msg
	}
}

func (l *tcpClient) log(msg string, err interface{}) {
	fmt.Println(msg, err)
}

func (l *tcpClient) LoopWrite() {
	for {
		select {
		case msg, ok := <-l.input:
			if !ok {
				return
			}

			bt := []byte(msg)
			_, err := l.conn.Write(bt)
			if err != nil {
				l.log("conn write", err)
				continue
			}
		case <-l.wait:
			return
		}
	}
}

func (l *tcpClient) LoopReceive() {
	for {
		select {
		case bt, ok := <-l.output:
			if !ok {
				return
			}
			l.handelReceive(bt)
		case <-l.wait:
			return
		}
	}
}

func (l *tcpClient) ReleaseChan() {
	select {
	case <-l.wait:
	default:
		close(l.output)
		close(l.input)
	}
}

func (l *tcpClient) Close() {
	l.conn.Close()
	l.ReleaseChan()
}

func (l *tcpClient) HasClosed() chan bool {
	return l.wait
}

func (l *tcpClient) Send(v string) {
	l.input <- v
}

func (l *tcpClient) OnReceive(f func(v []byte)) {
	l.receiveCallback = f
}

func NewTcpClient(addr string) *tcpClient {
	return &tcpClient{
		input:           make(chan string),
		output:          make(chan []byte),
		wait:            make(chan bool),
		conn:            nil,
		closeCallback:   nil,
		receiveCallback: nil,
		addr:            addr,
	}
}
