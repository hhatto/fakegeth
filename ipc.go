package main

import (
	"encoding/json"
	"log"
	"net"
	"os"
)

type IPCServer struct {
	path string
}

func NewIPCServer(path string) *IPCServer {
	return &IPCServer{path: path}
}

func (s *IPCServer) ListenAndServe() {
	endpoint := s.path
	os.Remove(endpoint)
	log.Printf("listen ipc, endpoint: %v", endpoint)
	l, err := net.Listen("unix", endpoint)
	if err != nil {
		log.Printf("uds:%s listen error: %v", endpoint, err)
		return
	}
	os.Chmod(endpoint, 0600)
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				log.Printf("ipc accept error: %v", err)
				return
			}
			defer conn.Close()

			//buf, err := bufio.NewReader(conn).ReadString('}')
			if err != nil {
				log.Printf("ipc read error: %v", err)
				return
			}

			decoder := json.NewDecoder(conn)
			var t RPCRequest
			if err := decoder.Decode(&t); err != nil {
				panic(err)
			}
			log.Printf("ipc: id=%v, method=%v, params=%v", t.Id, t.Method, t.Params)
		}
	}()
}

func (s *IPCServer) Shutdown() {
	os.Remove(s.path)
}
