package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
)

// Based off https://medium.com/ymedialabs-innovation/reverse-proxy-in-go-d26482acbcad
// and https://kasvith.me/posts/lets-create-a-simple-lb-go/

type server struct {
	target *url.URL
	proxy  *httputil.ReverseProxy
}

type serverPool struct {
	servers []*server
	current uint64
}

func (sp *serverPool) addServerToPool(server *server) {
	sp.servers = append(sp.servers, server)
}

var pool serverPool

func (f *Forego) startProxy(proxyPort int, backendPorts []int, of *OutletFactory) {

	for _, port := range backendPorts {
		u, _ := url.Parse(fmt.Sprintf("http://localhost:%d", port))
		p := httputil.NewSingleHostReverseProxy(u)
		pool.addServerToPool(&server{
			target: u,
			proxy:  p,
		})
	}

	reverseProxy := http.Server{
		Addr:    fmt.Sprintf(":%d", proxyPort),
		Handler: http.HandlerFunc(loadBalance),
	}

	if err := reverseProxy.ListenAndServe(); err != nil {
		of.SystemOutput(fmt.Sprintf("%v", err))
	}
}

func (sp *serverPool) getServer() *server {
	return sp.servers[(atomic.AddUint64(&sp.current, uint64(1)) % uint64(len(sp.servers)))]
}

func loadBalance(w http.ResponseWriter, r *http.Request) {
	s := pool.getServer()
	if s != nil {
		s.proxy.ServeHTTP(w, r)
		return
	}
	http.Error(w, "Service not available", http.StatusServiceUnavailable)
}
