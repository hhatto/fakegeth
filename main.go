package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"

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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var t RPCRequest
	if err := decoder.Decode(&t); err != nil {
		panic(err)
	}

	log.Printf("method: %v, user-agent: %v, body: %v", r.Method, r.Header["User-Agent"], t)
	resJson := RPCResponse{
		Id:      t.Id,
		Jsonrpc: "2.0",
		Result:  "ok",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resJson)
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("websocket upgrade: %v", err)
		return
	}
	defer c.Close()

	for {
		mt, msg, err := c.ReadMessage()
		if err != nil {
			log.Printf("read error: %v", err)
			break
		}
		log.Printf("recvmsg: %v", string(msg))

		err = c.WriteMessage(mt, msg)
		if err != nil {
			log.Printf("write error: %v", err)
			break
		}
	}
}

func main() {
	// signal handle
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// load config
	var config *Config
	var err error
	if len(os.Args) >= 2 {
		config, err = LoadConfig(os.Args[1])
		if err != nil {
			log.Printf("read config error: %v", err)
			return
		}
	}
	log.Printf("config: %v", config)

	// listen http
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	httpAddr := ":8545"
	log.Printf("listen http, endpoint: %v", httpAddr)
	go http.ListenAndServe(httpAddr, mux)

	// listen ipc
	endpoint := "./fakegeth.ipc"
	os.Remove(endpoint)
	log.Printf("listen ipc, endpoint: %v", endpoint)
	l, err := net.Listen("unix", endpoint)
	if err != nil {
		log.Printf("uds:%s listen error: %v", endpoint, err)
		return
	}
	os.Chmod(endpoint, 0600)
	defer os.Remove(endpoint)
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				log.Printf("ipc accept error: %v", err)
				return
			}
			defer conn.Close()

			buf, err := bufio.NewReader(conn).ReadString('}')
			if err != nil {
				log.Printf("ipc read error: %v", err)
				return
			}

			log.Printf("ipc read, %d bytes. recvmsg: %v", len(buf), string(buf))
		}
	}()

	// listen websocket
	wsAddr := ":8546"
	log.Printf("listen websocket, endpoint: %v", wsAddr)
	wsMux := http.NewServeMux()
	wsMux.HandleFunc("/", websocketHandler)
	go http.ListenAndServe(wsAddr, wsMux)

	// wait signal
	_ = <-c
	log.Println("stop")
}
