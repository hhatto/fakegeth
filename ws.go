package main

import (
	"log"
	"net/http"
)

type WebsocketServer struct {
	mux *http.ServeMux
}

func NewWebsocketServer() *WebsocketServer {
	s := &WebsocketServer{
		mux: http.NewServeMux(),
	}
	s.mux.HandleFunc("/", s.wsHandler)
	return s
}

func (s *WebsocketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *WebsocketServer) wsHandler(w http.ResponseWriter, r *http.Request) {
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
