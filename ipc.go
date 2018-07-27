package main

import (
	"bufio"
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

			buf, err := bufio.NewReader(conn).ReadString('}')
			if err != nil {
				log.Printf("ipc read error: %v", err)
				return
			}

			log.Printf("ipc read, %d bytes. recvmsg: %v", len(buf), string(buf))
		}
	}()
}

func (s *IPCServer) Shutdown() {
	os.Remove(s.path)
}
