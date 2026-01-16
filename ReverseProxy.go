package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"net"
	"time"
)

var (
	counter  int
	mu       sync.Mutex
	backends = []*url.URL{
		mustParse("http://localhost:8081"),
		mustParse("http://localhost:8082"),
		mustParse("http://localhost:8083"),
	}
	rateLimiter = make(map[string]int)
	limitePerMinute = 10
)

func mustParse(raw string) *url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		log.Fatal(err)
	}
	return u
}

func resetRateLimiter() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		mu.Lock()
		rateLimiter = make(map[string]int)
		mu.Unlock()

		fmt.Println("Rate limiter reset")
	}
}


func handler(w http.ResponseWriter, r *http.Request) {

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		// fallback (very rare, but don't be stupid)
		ip = r.RemoteAddr
	}
	//Figured there would a datarace here since http requests are sent in parallel
	mu.Lock()
	backend := backends[counter%len(backends)] //This line distributes http requests over all backend servers evenly
	counter++
	rateLimiter[ip]++
	mu.Unlock()

	
	if rateLimiter[ip] < limitePerMinute {
		proxy := httputil.NewSingleHostReverseProxy(backend)
		w.Write([]byte("Your request was forwarded to Server "+ fmt.Sprintf("%d",counter%len(backends)+ 1) + "\n" ))
		proxy.ServeHTTP(w, r)
		fmt.Println(r.Header)
	}

	}

func main() {
	go resetRateLimiter()
	http.HandleFunc("/", handler)
	http.ListenAndServe(":9090", nil)
}

