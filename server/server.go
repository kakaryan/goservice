package server

import (
	"fmt"
	"github.com/kakaryan/goservice/service"
	"github.com/rcrowley/goagain"
	"log"
	"net"
	"time"
)

type Server struct {
	ip      string
	port    string
	handler map[string]service.Handler
}

func NewServer(ip string, port string, h map[string]service.Handler) *Server {
	s := &Server{
		ip:      ip,
		port:    port,
		handler: h,
	}
	return s
}

func (s *Server) RunServer() {
	l, err := goagain.Listener()

	if nil != err {

		// Listen on a TCP or a UNIX domain socket (TCP here).
		// l, err = net.Listen("tcp", "127.0.0.1:48879")
		addr := fmt.Sprintf("%s:%s", s.ip, s.port)
		l, err = net.Listen("tcp", addr)
		if nil != err {
			log.Fatalln(err)
		}
		log.Println("listening on", l.Addr())

		// Accept connections in a new goroutine.
		s.WaitForConnections(l)

	} else {
		log.Println("resuming listening on", l.Addr())
		s.WaitForConnections(l)

		// Kill the parent, now that the child has started successfully.
		if err := goagain.Kill(); nil != err {
			log.Fatalln(err)
		}
	}

	// Block the main goroutine awaiting signals.
	if _, err := goagain.Wait(l); nil != err {
		log.Fatalln(err)
	}
	if err := l.Close(); nil != err {
		log.Fatalln(err)
	}
	time.Sleep(1e9)
}

func (s *Server) WaitForConnections(ls net.Listener) {
	reqChannel := make(chan service.ChanReq)
	for {
		conn, e := ls.Accept()
		if e == nil {
			host, port, err := net.SplitHostPort(conn.RemoteAddr().String())
			if nil != err {
				log.Fatalln(err)
			}
			go service.RunService(reqChannel, s.handler)
			handler := &service.ReqHandler{reqChannel}
			log.Println("Remote Addr", host, "Remote Port", port)
			go Handle(conn, handler)
		} else {
			log.Println("Error accepting from %s", ls)
		}
	}
}
