package main

import (
	"fmt"
	"github.com/kakaryan/goservice/server"
	"github.com/kakaryan/goservice/service"
	"log"
	"math/rand"
	"time"
	// "net"
	// "os"
	// "os/signal"
	"syscall"
	// "time"
)

func init() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
	log.SetPrefix(fmt.Sprintf("pid:%d ", syscall.Getpid()))
}

func main() {
	var handlers = map[string]service.Handler{
		// "test":             handleTest,
		// "ping": handlePing,
		"RPC": handleRpc,
	}
	s := server.NewServer("127.0.0.1", "48879", handlers)
	s.RunServer()
}

func handleRpc(req *service.Request) *service.Response {
	var response service.Response
	t := rand.Intn(10)
	nums := []int{21, 22, 23, 24, 25, 26, 27, 28, 29, 30}
	sleep := time.Duration(nums[t]) * time.Second
	res := fmt.Sprintf("sleep for %d second", sleep/time.Second)
	response.Data = []byte(res)
	log.Printf("response is %s", response.Data)
	time.Sleep(sleep)
	return &response
}
