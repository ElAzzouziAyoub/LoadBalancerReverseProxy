package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

var (
	counter  int
	mu       sync.Mutex
	backends = []*url.URL{
		mustParse("http://localhost:8081"),
		mustParse("http://localhost:8082"),
		mustParse("http://localhost:8083"),
	}
)

func mustParse(raw string) *url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		panic(err)
	}
	return u
}

func handler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	backend := backends[counter%len(backends)]
	counter++
	mu.Unlock()

	proxy := httputil.NewSingleHostReverseProxy(backend)
	proxy.ServeHTTP(w, r)
}

func main() {
	http.ListenAndServe(":9090", http.HandlerFunc(handler))
}

