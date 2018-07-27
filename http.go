package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type HTTPServer struct {
	mux *http.ServeMux
}

func NewHTTPServer() *HTTPServer {
	s := &HTTPServer{
		mux: http.NewServeMux(),
	}
	s.mux.HandleFunc("/", s.indexHandler)
	return s
}

func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *HTTPServer) indexHandler(w http.ResponseWriter, r *http.Request) {
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
