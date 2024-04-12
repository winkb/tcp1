package myws

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/gorilla/websocket"
	"github.com/winkb/tcp1/btmsg"
	. "github.com/winkb/tcp1/contracts"
	"github.com/winkb/tcp1/util"
)

type wrapConn struct {
	*websocket.Conn
}

func (l *wrapConn) Read(b []byte) (n int, err error) {
	panic("implement me")
}

func (l *wrapConn) SetDeadline(t time.Time) error {
	if err := l.Conn.SetReadDeadline(t); err != nil {
		return err
	}
	return l.Conn.SetWriteDeadline(t)
}

func (l *wrapConn) GetRemoteIp() string {
	return l.Conn.RemoteAddr().String()
}

// Write writes data to the connection.
// Write can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetWriteDeadline.
func (l *wrapConn) Write(b []byte) (n int, err error) {
	return
}

var upgrader = websocket.Upgrader{} // use default options

var _ ITcpServer = (*Ws)(nil)

func NewWs(addr string, wsPath string, r btmsg.IMsgReader) *Ws {
	return &Ws{
		wsPath:          wsPath,
		addr:            addr,
		reader:          r,
		closeCallback:   func(conn *TcpConn, isServer, isClient bool) {},
		receiveCallback: func(conn *TcpConn, msg btmsg.IMsg) {},
		conns:           sync.Map{},
		lastId:          0,
		stop:            0,
		lock:            sync.RWMutex{},
		timeout:         time.Second * 3,
	}
}

type Ws struct {
	wsPath          string
	listener        *websocket.Conn
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

func (l *Ws) Shutdown() {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.stop != 0 {
		return
	}

	l.stop = 2

	l.conns.Range(func(key, value any) bool {
		v, ok := value.(*TcpConn)
		if ok {
			_ = v.Conn.Close()
		}
		return true
	})

	err := l.listener.Close()
	if err != nil {
		fmt.Println(err)
	}
}
func (l *Ws) Send(conn *TcpConn, v btmsg.IMsg) {
	l.lock.RLock()
	if l.stop != 0 {
		l.lock.RUnlock()
		return
	}
	l.lock.RUnlock()

	conn.Lock.RLock()
	if conn.IsClose {
		conn.Lock.RUnlock()
		return
	}
	conn.Lock.RUnlock()

	conn.Input <- v
}

func (l *Ws) getConnById(id uint32) (conn *TcpConn, ok bool) {
	v, o := l.conns.Load(id)
	if !o {
		return
	}

	conn, ok = v.(*TcpConn)

	return
}

func (l *Ws) SendById(id uint32, v btmsg.IMsg) {
	conn, ok := l.getConnById(id)
	if !ok {
		log.Err(errors.Errorf("not found conn %d", id))
		return
	}

	l.Send(conn, v)
}

func (l *Ws) OnReceive(f ServerReceiveCallback) {
	l.receiveCallback = f
}

func (l *Ws) OnClose(f ServerCloseCallback) {
	l.closeCallback = f
}

func (l *Ws) listen() (err error) {
	return nil
}

func (l *Ws) getConnAutoIncId() uint32 {
	for {
		val := atomic.LoadUint32(&l.lastId)
		old := val
		val += 1
		if atomic.CompareAndSwapUint32(&l.lastId, old, val) {
			return val
		}
	}
}

func (l *Ws) Start() (wg *sync.WaitGroup, err error) {
	wg = &sync.WaitGroup{}
	// conn server
	err = l.listen()
	if err != nil {
		return
	}
	// read
	util.MyGoWg(wg, "conn_accept", func() {
		l.LoopAccept(func(conn *websocket.Conn) {
			// 注意 这里不能阻塞 lock,因为accept，有lock判断

			newId := l.getConnAutoIncId()
			myConn := &TcpConn{
				Conn: &wrapConn{
					Conn:conn,
				},
				Id:       newId,
				Input:    make(chan btmsg.IMsg),
				Output:   make(chan btmsg.IMsg),
				WaitConn: make(chan bool),
			}

			util.MyGoWg(wg, fmt.Sprintf("%d_conn_read", newId), func() {
				l.LoopRead(myConn)
			})

			util.MyGoWg(wg, fmt.Sprintf("%d_conn_consume_input", newId), func() {
				l.ConsumeInput(myConn)
			})

			util.MyGoWg(wg, fmt.Sprintf("%d_conn_consume_output", newId), func() {
				l.ConsumeOutput(myConn)
			})

			fmt.Println(conn.RemoteAddr().String() + "conn success")

			l.saveConn(newId, myConn)
		})
	})

	fmt.Println("start server " + l.addr)

	return
}

func (l *Ws) saveConn(id uint32, conn *TcpConn) {
	l.conns.Store(id, conn)
}

func (l *Ws) LoopAccept(f func(conn *websocket.Conn)) {
	http.HandleFunc("/"+l.wsPath, func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}

		f(conn)
	})

	err := http.ListenAndServe(l.addr, nil)
	if err != nil {
		panic(err)
	}
}

func (l *Ws) ConsumeInput(conn *TcpConn) {
	for {
		select {
		case <-conn.WaitConn:
			return
		case msg := <-conn.Output:
			l.handelReceive(conn, msg)
		}
	}
}

func (l *Ws) handelReceive(conn *TcpConn, bt btmsg.IMsg) {
	if l.receiveCallback != nil {
		l.receiveCallback(conn, bt)
	}
}

func (l *Ws) ConsumeOutput(conn *TcpConn) {
	for {
		select {
		case <-conn.WaitConn:
			return
		case msg := <-conn.Output:
			l.handelReceive(conn, msg)
		}
	}
}

func (l *Ws) LoopRead(conn *TcpConn) {
	defer func() {
		select {
		case <-conn.WaitConn:
		default:
			close(conn.WaitConn)
		}
	}()
	for {
		select {
		case <-conn.WaitConn:
			return
		default:
			res := l.reader.ReadMsg(conn.Conn)
			err := res.GetErr()
			conn.Lock.Lock()
			if err != nil {
				conn.IsClose = true
				l.removeConn(conn.Id)
			}
			conn.Lock.Unlock()

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

			conn.Output <- res.GetMsg()
		}
	}
}

func (l *Ws) Broadcast(bt btmsg.IMsg) {
	l.conns.Range(func(key, value any) bool {
		v, ok := value.(*TcpConn)
		if !ok {
			return false
		}
		l.Send(v, bt)
		return true
	})
}

func (l *Ws) removeConn(id uint32) {
	l.conns.Delete(id)
}

func (l *Ws) handelReadClose(conn *TcpConn, isServer bool, isClient bool) {
	close(conn.WaitConn)
	if l.closeCallback != nil {
		l.closeCallback(conn, isServer, isClient)
	}
}
