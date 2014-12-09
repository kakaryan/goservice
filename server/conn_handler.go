package server

import (
	service "github.com/kakaryan/goservice/service"
	"io"
	"log"
)

type Transmitter interface {
	Encode()
}

type RequestHandler interface {
	HandleMessage(io.Writer, *service.Request) *service.Response
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func Handle(s io.ReadWriteCloser, h RequestHandler) error {
	defer func() { must(s.Close()) }()
	var err error
	for err == nil {
		err = HandleMessage(s, s, h)
	}
	return err
}

func ReadPacket(r io.Reader) (rv service.Request, err error) {
	rv.Receive(r)
	return
}

// Handle an individual message.
func HandleMessage(r io.Reader, w io.Writer, handler RequestHandler) error {
	req, err := ReadPacket(r)
	log.Printf("%s", req)
	if err != nil {
		return err
	}
	res := handler.HandleMessage(w, &req)
	log.Printf("Got a response: %s", res)
	if res == nil {
		// Quiet command
		return nil
	}
	_, err = res.Transmit(w)
	if err != nil {
		return err
	}

	return io.EOF
}
