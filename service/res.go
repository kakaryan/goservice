package service

import (
	"bufio"
	"encoding/json"
	"io"
	"strconv"
)

// A Response
type Response struct {
	// body
	Data []byte
}

// Send this response message across a writer.
func (res *Response) Transmit(w io.Writer) (n int, err error) {
	data := res.Bytes()
	l := strconv.Itoa(len(data))
	l += "\n"
	resWriter := bufio.NewWriter(w)
	n, err = resWriter.WriteString(l)
	n, err = resWriter.WriteString(string(res.Bytes()) + "\n")
	resWriter.Flush()
	return
}

// The actual bytes transmitted for this response.
func (res *Response) Bytes() []byte {
	data, _ := json.Marshal(string(res.Data))
	return data
}
