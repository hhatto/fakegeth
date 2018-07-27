package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type RPCRequest struct {
	Jsonrpc string
	Method  string
	Params  []interface{}
	Id      int
}

type JSONRPCError struct {
	Code    int
	Message string
}

type RPCResponse struct {
	Id      int
	Jsonrpc string
	Result  string `json:"result,omitempty"`
	Error   JSONRPCError
}

var upgrader = websocket.Upgrader{}

func serveHTTP(config *Config, wg *sync.WaitGroup, stop <-chan struct{}) {
	addr := fmt.Sprintf("%s:%d", config.Http.Host, config.Http.Port)
	log.Printf("listen http, endpoint: %v", addr)
	ss := NewHTTPServer()
	s := &http.Server{Addr: addr, Handler: ss}
	go s.ListenAndServe()

	go func() {
		wg.Add(1)
		<-stop
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		s.Shutdown(ctx)
		log.Println("stop http")
		wg.Done()
	}()
}

func serveWebsocket(config *Config, wg *sync.WaitGroup, stop <-chan struct{}) {
	addr := fmt.Sprintf("%s:%d", config.Websocket.Host, config.Websocket.Port)
	log.Printf("listen websocket, endpoint: %v", addr)
	ss := NewWebsocketServer()
	s := &http.Server{Addr: addr, Handler: ss}
	go s.ListenAndServe()

	go func() {
		wg.Add(1)
		<-stop
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		s.Shutdown(ctx)
		log.Println("stop websocket")
		wg.Done()
	}()
}

func serveIPC(config *Config, wg *sync.WaitGroup, stop <-chan struct{}) {
	endpoint := config.Ipc.Path
	s := NewIPCServer(endpoint)
	go s.ListenAndServe()

	go func() {
		wg.Add(1)
		<-stop
		s.Shutdown()
		log.Println("stop ipc")
		wg.Done()
	}()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: fakegeth CONFFILE")
		return
	}

	// signal handle
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// load config
	config, err := LoadConfig(os.Args[1])
	if err != nil {
		log.Printf("read config error: %v", err)
		return
	}
	log.Printf("config: %v, %v, %v", config.Http, config.Websocket, config.Ipc)

	stop := make(chan struct{})
	var wg sync.WaitGroup
	// listen http
	if config != nil && config.Http != nil {
		go serveHTTP(config, &wg, stop)
	}

	// listen ipc
	if config != nil && config.Ipc != nil {
		go serveIPC(config, &wg, stop)
	}

	// listen websocket
	if config != nil && config.Websocket != nil {
		serveWebsocket(config, &wg, stop)
	}

	// wait signal
	_ = <-c

	close(stop)
	wg.Wait()
	log.Println("stop")
}
