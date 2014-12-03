package service

import (
	"io"
	// "log"
)

type ChanReq struct {
	Req *Request
	Res chan *Response
}

type ReqHandler struct {
	Ch chan ChanReq
}

func (rh *ReqHandler) HandleMessage(w io.Writer, req *Request) *Response {
	cr := ChanReq{
		req,
		make(chan *Response),
	}
	rh.Ch <- cr
	return <-cr.Res
}
