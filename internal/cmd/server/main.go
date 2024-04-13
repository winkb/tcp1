package main

import (
	"fmt"
	"github.com/winkb/tcp1/internal/cmd/server/handles"
	"github.com/winkb/tcp1/net/myws"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/winkb/tcp1/btmsg"
	"github.com/winkb/tcp1/contracts"
)

func main() {
	//var server = mytcp.NewTcpServer("989", btmsg.NewReader(func() btmsg.IHead {
	//	return btmsg.NewMsgHeadTcp()
	//}))

	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		err := homeTemplate.Execute(w, "ws://"+r.Host+"/ws")
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
		}
	})

	var server = myws.NewWs("localhost:9899", "ws", btmsg.NewReader(func() btmsg.IHead {
		return btmsg.NewMsgHeadWs()
	}))
	wg, err := server.Start()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", server.LoopAccept)

	go func() {
		wg.Add(1)
		defer wg.Done()

		fmt.Println("localhost:9899")

		http.ListenAndServe("localhost:9899", nil)
	}()

	server.OnClose(func(s contracts.ITcpServer, conn *contracts.TcpConn, isServer bool, isClient bool) {
		if isClient {
			fmt.Println("客户端断开连接")
		}

		if isServer {
			fmt.Println("我自己断开连接")
		}
	})

	server.OnReceive(func(s contracts.ITcpServer, conn *contracts.TcpConn, msg btmsg.IMsg) {
		act := msg.GetAct()
		hv, ok := handles.Routes[act]
		if !ok {
			fmt.Println("not found handle", act)

			// 走默认路由
			act = 0
			hv = handles.Routes[act]
		}

		hv.Handle(s, conn, msg)
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

var homeTemplate = template.Must(template.New("").Parse(`<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {

    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;

    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
        output.scroll(0, output.scrollHeight);
    };

    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };

    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
		let msg  = {
			act:1, 
			data:{
				content:input.value
			}
		};
        ws.send(JSON.stringify(msg));
        return false;
    };

    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };

});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button type="button" id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output" style="max-height: 70vh;overflow-y: scroll;"></div>
</td></tr></table>
</body>
</html>`))
