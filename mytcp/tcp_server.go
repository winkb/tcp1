package mytcp

import (
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/winkb/tcp1/btmsg"
)

type ServerCloseCallback func(conn *TcpConn, isServer bool, isClient bool)
type ServerReceiveCallback func(conn *TcpConn, msg btmsg.IMsg)

type ITcpServer interface {
	Shutdown()
	Send(conn *TcpConn, v btmsg.IMsg)
	SendById(id uint32, v btmsg.IMsg)
	OnReceive(f ServerReceiveCallback)
	OnClose(f ServerCloseCallback)
	Start() (wg *sync.WaitGroup, err error)
	LoopAccept(f func(conn net.Conn))
	ConsumeInput(conn *TcpConn)
	ConsumeOutput(conn *TcpConn)
	LoopRead(conn *TcpConn)
	Broadcast(bt btmsg.IMsg)
}

var _ ITcpServer = (*tcpServer)(nil)

type tcpServer struct {
	listener        net.Listener
	closeCallback   ServerCloseCallback
	receiveCallback ServerReceiveCallback
	addr            string
	conns           sync.Map
	lastId          uint32
	stop            int
	lock            sync.RWMutex
	reader          btmsg.IMsgReader
	timeout         time.Duration
}

type TcpConn struct {
	conn     net.Conn
	id       uint32
	input    chan btmsg.IMsg
	output   chan btmsg.IMsg
	waitConn chan bool
	lock     sync.RWMutex
	isClose  bool
}

func (l *TcpConn) GetRemoteIp() string {
	if l.conn == nil {
		return ""
	}
	return l.conn.RemoteAddr().String()
}

func (l *TcpConn) GetId() uint32 {
	return l.id
}

func NewTcpServer(port string, r btmsg.IMsgReader) *tcpServer {
	return &tcpServer{
		listener: nil,
		closeCallback: func(conn *TcpConn, isServer bool, isClient bool) {
		},
		receiveCallback: func(conn *TcpConn, msg btmsg.IMsg) {
		},
		addr:    ":" + port,
		conns:   sync.Map{},
		lastId:  0,
		stop:    0,
		lock:    sync.RWMutex{},
		reader:  r,
		timeout: time.Second * 3,
	}
}

func (l *tcpServer) LoopAccept(f func(conn net.Conn)) {
	for {
		accept, err := l.listener.Accept()
		if err != nil {
			if _, ok := err.(*net.OpError); ok {
				fmt.Println("server shutdown")
				return
			}

			log.Err(errors.Wrap(err, "accept"))
			return
		}

		l.lock.RLock()
		if l.stop != 0 {
			fmt.Println("server is stop")
			l.lock.RUnlock()
			continue
		}

		f(accept)
		l.lock.RUnlock()
	}
}

func (l *tcpServer) getConnAutoIncId() uint32 {
	for {
		val := atomic.LoadUint32(&l.lastId)
		old := val
		val += 1
		if atomic.CompareAndSwapUint32(&l.lastId, old, val) {
			return val
		}
	}
}

func (l *tcpServer) getConnById(id uint32) (conn *TcpConn, ok bool) {
	v, o := l.conns.Load(id)
	if !o {
		return
	}

	conn, ok = v.(*TcpConn)

	return
}

func (l *tcpServer) saveConn(id uint32, conn *TcpConn) {
	l.conns.Store(id, conn)
}

func (l *tcpServer) removeConn(id uint32) {
	l.conns.Delete(id)
}

func (l *tcpServer) ConsumeOutput(conn *TcpConn) {
	for {
		select {
		case <-conn.waitConn:
			return
		case msg := <-conn.output:
			l.handelReceive(conn, msg)
		}
	}
}

func (l *tcpServer) writeSend(conn *TcpConn, msg btmsg.IMsg) {
	l.lock.RLock()
	defer l.lock.Unlock()

	if l.stop != 0 {
		return
	}

	_ = conn.conn.SetWriteDeadline(time.Now().Add(l.timeout))
	id := conn.id

	var err error
	conn.lock.RLock()
	defer conn.lock.RUnlock()
	if conn.isClose {
		log.Print("conn is closed, drop msg")
		return
	}

	_, err = conn.conn.Write(msg.ToSendByte())
	if err != nil {
		log.Err(errors.Wrapf(err, "conn %d write err", id))
		return
	}

	log.Print("input id", id, "msg", msg.GetAct(), string(msg.BodyByte()))
}

func (l *tcpServer) ConsumeInput(conn *TcpConn) {
	for {
		select {
		case <-conn.waitConn:
			return
		case msg := <-conn.input:
			l.writeSend(conn, msg)
		}
	}
}

func (l *tcpServer) LoopRead(conn *TcpConn) {
	defer func() {
		select {
		case <-conn.waitConn:
		default:
			close(conn.waitConn)
		}
	}()
	for {
		select {
		case <-conn.waitConn:
			return
		default:
			res := l.reader.ReadMsg(conn.conn)
			err := res.GetErr()
			conn.lock.Lock()
			if err != nil {
				conn.isClose = true
				l.removeConn(conn.id)
			}
			conn.lock.Unlock()

			if err != nil {

				if res.IsCloseByClient() {
					l.handelReadClose(conn, false, true)
					return
				}

				if res.IsCloseByServer() {
					l.handelReadClose(conn, true, true)
					return
				}

				log.Err(errors.Wrap(err, "read"))
				return
			}

			conn.output <- res.GetMsg()
		}
	}
}

func (l *tcpServer) handelReadClose(conn *TcpConn, isServer bool, isClient bool) {
	close(conn.waitConn)
	if l.closeCallback != nil {
		l.closeCallback(conn, isServer, isClient)
	}
}

func (l *tcpServer) handelReceive(conn *TcpConn, bt btmsg.IMsg) {
	if l.receiveCallback != nil {
		l.receiveCallback(conn, bt)
	}
}

func (l *tcpServer) Shutdown() {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.stop != 0 {
		return
	}

	l.stop = 2

	l.conns.Range(func(key, value any) bool {
		v, ok := value.(*TcpConn)
		if ok {
			_ = v.conn.Close()
		}
		return true
	})

	err := l.listener.Close()
	if err != nil {
		fmt.Println(err)
	}
}

func (l *tcpServer) Send(conn *TcpConn, v btmsg.IMsg) {
	l.lock.RLock()
	if l.stop != 0 {
		l.lock.RUnlock()
		return
	}
	l.lock.RUnlock()

	conn.lock.RLock()
	if conn.isClose {
		conn.lock.RUnlock()
		return
	}
	conn.lock.RUnlock()

	conn.input <- v
}

func (l *tcpServer) SendById(id uint32, v btmsg.IMsg) {
	conn, ok := l.getConnById(id)
	if !ok {
		log.Err(errors.Errorf("not found conn %d", id))
		return
	}

	l.Send(conn, v)
}

func (l *tcpServer) OnReceive(f ServerReceiveCallback) {
	l.receiveCallback = f
}

func (l *tcpServer) OnClose(f ServerCloseCallback) {
	l.closeCallback = f
}

func (l *tcpServer) Start() (wg *sync.WaitGroup, err error) {
	wg = &sync.WaitGroup{}
	// conn server
	err = l.listen()
	if err != nil {
		return
	}
	// read
	MyGoWg(wg, "conn_accept", func() {
		l.LoopAccept(func(conn net.Conn) {
			// 注意 这里不能阻塞 lock,因为accept，有lock判断

			newId := l.getConnAutoIncId()
			myConn := &TcpConn{
				conn:     conn,
				id:       newId,
				input:    make(chan btmsg.IMsg),
				output:   make(chan btmsg.IMsg),
				waitConn: make(chan bool),
			}

			MyGoWg(wg, fmt.Sprintf("%d_conn_read", newId), func() {
				l.LoopRead(myConn)
			})

			MyGoWg(wg, fmt.Sprintf("%d_conn_consume_input", newId), func() {
				l.ConsumeInput(myConn)
			})

			MyGoWg(wg, fmt.Sprintf("%d_conn_consume_output", newId), func() {
				l.ConsumeOutput(myConn)
			})

			fmt.Println(conn.RemoteAddr().String() + "conn success")

			l.saveConn(newId, myConn)
		})
	})

	fmt.Println("start server " + l.addr)

	return
}

func (l *tcpServer) listen() (err error) {
	var conn net.Listener
	conn, err = net.Listen("tcp", l.addr)
	if err != nil {
		err = errors.Wrap(err, "dial:"+l.addr)
		return
	}
	l.listener = conn
	return
}

func (l *tcpServer) Broadcast(bt btmsg.IMsg) {
	l.conns.Range(func(key, value any) bool {
		v, ok := value.(*TcpConn)
		if !ok {
			return false
		}
		l.Send(v, bt)
		return true
	})
}
