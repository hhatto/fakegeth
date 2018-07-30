package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
)

func serveProxy() {
	var addr string = "localhost:9545"
	log.Println("Starting proxy server on", addr)

	director := func(req *http.Request) {
		url := *req.URL
		url.Scheme = "http"
		url.Host = "127.0.0.1:8545"

		buf := new(bytes.Buffer)
		if req.Body != nil {
			buff := new(bytes.Buffer)
			io.Copy(buff, req.Body)
			s := buff.String()
			buf.WriteString(s)

			decoder := json.NewDecoder(buff)
			var t RPCRequest
			if err := decoder.Decode(&t); err != nil {
				panic(err)
			}

			log.Printf("proxy: user-agent=%v, id=%v, method=%v, params=%v",
				req.Header["User-Agent"], t.Id, t.Method, t.Params)
		}

		proxyReq, err := http.NewRequest(req.Method, url.String(), buf)
		if err != nil {
			log.Fatal(err.Error())
		}
		proxyReq.Header = req.Header
		*req = *proxyReq
	}

	rp := &httputil.ReverseProxy{Director: director}
	server := http.Server{
		Addr:    addr,
		Handler: rp,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
