package mytcp

import (
	"fmt"
	"github.com/winkb/tcp1/btmsg"
	"github.com/winkb/tcp1/util"
	"net"
	"sync"

)

type clientReceiveCallback func(msg btmsg.IMsg)
type clientCloseCallback func(isServer bool, isClient bool)

type ITcpClient interface {
	LoopRead()
	ReleaseChan()
	LoopWrite()
	LoopReceive()
	Close()
	Send(v btmsg.IMsg)
	OnReceive(f clientReceiveCallback)
	OnClose(f clientCloseCallback)
	Start() (wg *sync.WaitGroup, err error)
	HasClosed() chan bool
}

var _ ITcpClient = (*tcpClient)(nil)

type tcpClient struct {
	input           chan btmsg.IMsg
	output          chan btmsg.IMsg
	wait            chan bool
	conn            net.Conn
	closeCallback   clientCloseCallback
	receiveCallback clientReceiveCallback
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
	util.MyGoWg(wg, "conn_read", l.LoopRead)
	// write
	util.MyGoWg(wg, "conn_write", l.LoopWrite)
	// on msg
	util.MyGoWg(wg, "conn_receive", l.LoopReceive)

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

func (l *tcpClient) handelReceive(msg btmsg.IMsg) {
	if l.receiveCallback != nil {
		l.receiveCallback(msg)
	}
}

func (l *tcpClient) OnClose(f clientCloseCallback) {
	l.closeCallback = f
}

func (l *tcpClient) LoopRead() {

	for {
		r := btmsg.NewReader()
		res := r.ReadMsg(l.conn)
		if err := res.GetErr(); err != nil {
			if res.IsCloseByServer() {
				l.handelReadClose(true, false)
				return
			}

			if res.IsCloseByClient() {
				l.handelReadClose(false, true)
				return
			}

			if err != nil {
				l.log("conn read", err)
				continue
			}
		}

		l.output <- res.GetMsg()
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

			_, err := l.conn.Write(msg.ToSendByte())
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

func (l *tcpClient) Send(v btmsg.IMsg) {
	l.input <- v
}

func (l *tcpClient) OnReceive(f clientReceiveCallback) {
	l.receiveCallback = f
}

func NewTcpClient(addr string) *tcpClient {
	return &tcpClient{
		input:           make(chan btmsg.IMsg),
		output:          make(chan btmsg.IMsg),
		wait:            make(chan bool),
		conn:            nil,
		closeCallback:   nil,
		receiveCallback: nil,
		addr:            addr,
	}
}
