package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:48879")

	defer conn.Close()

	if err != nil {
		log.Panicln("Connection error")
	}
	w := bufio.NewWriter(conn)
	_, err = w.WriteString("3\n")
	_, err = w.WriteString("RPC\n")
	_, err = w.WriteString("4\n")
	_, err = w.WriteString("asdf\n")
	w.Flush()
	s := bufio.NewScanner(conn)
	for s.Scan() {
		log.Println(s.Text())
	}
	if err := s.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
