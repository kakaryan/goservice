package main

import (
	"fmt"
	"github.com/kakaryan/goservice/server"
	"github.com/kakaryan/goservice/service"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
	log.SetPrefix(fmt.Sprintf("pid:%d ", syscall.Getpid()))
}

func main() {

	// Listen on a TCP or a UNIX domain socket (TCP here).
	l, err := net.Listen("tcp", "127.0.0.1:48879")
	if nil != err {
		log.Fatalln(err)
	}
	log.Println("listening on", l.Addr())

	waitForConnections(l)

	// 阻塞主协程等待信号.
	if _, err := Wait(l); nil != err {
		log.Fatalln(err)
	}
	if err := l.Close(); nil != err {
		log.Fatalln(err)
	}
	time.Sleep(1e9)

}

func Wait(l net.Listener) (syscall.Signal, error) {
	ch := make(chan os.Signal, 2)
	signal.Notify(
		ch,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)
	for {
		sig := <-ch
		log.Println(sig.String())
		switch sig {

		// SIGINT should exit.
		case syscall.SIGINT:
			return syscall.SIGINT, nil

		// SIGQUIT should exit gracefully.
		case syscall.SIGQUIT:
			return syscall.SIGQUIT, nil

		// SIGTERM should exit.
		case syscall.SIGTERM:
			return syscall.SIGTERM, nil

		}
	}
}

func waitForConnections(ls net.Listener) {
	reqChannel := make(chan service.ChanReq)
	for {
		conn, e := ls.Accept()
		if e == nil {
			host, port, err := net.SplitHostPort(conn.RemoteAddr().String())
			if nil != err {
				log.Fatalln(err)
			}
			go service.RunService(reqChannel)
			handler := &service.ReqHandler{reqChannel}
			log.Println("Remote Addr", host, "Remote Port", port)
			go server.Handle(conn, handler)
		} else {
			log.Println("Error accepting from %s", ls)
		}
	}
}
