package service

import (
	"bufio"
	"errors"
	"io"
	"log"
	"strconv"
)

type Request struct {
	Opcode []byte
	Data   []byte
}

func (req *Request) Receive(r io.Reader) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in Receive", r)
		}
	}()
	buf := make([]byte, 1024)
	input := bufio.NewScanner(r)
	hl, bl, i := 0, 0, 0
	for input.Scan() {
		buf = input.Bytes()
		switch i {
		case 0:
			hl, _ = strconv.Atoi(string(buf))
		case 1:
			req.Opcode = buf
			if len(buf) != hl {
				log.Panicln(errors.New("Header length mismatched."))
			}
		case 2:
			bl, _ = strconv.Atoi(string(buf))
		case 3:
			req.Data = buf
			if len(buf) != bl {
				log.Panicln(errors.New("Body length mismatched."))
			}
			return
		default:
			log.Panicln(errors.New("Data content mismatched."))
		}
		i++
	}
	if 0 == i {
		return
	}
	if 4 != i {
		log.Println("field: %s != 3", i)
		panic(errors.New("Data field mismatched."))
	}
}
