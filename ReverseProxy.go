package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"
)

type Backend struct {
	URL         *url.URL
	Alive       bool
	Connections int
}

var (
	counter int
	mu      sync.Mutex

	backends = []*Backend{
		{URL: mustParse("http://localhost:8081"), Alive: true},
		{URL: mustParse("http://localhost:8082"), Alive: true},
		{URL: mustParse("http://localhost:8083"), Alive: true},
	}

	rateLimiter      = make(map[string]int)
	limitePerMinute  = 10
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


func getNextBackend() *Backend {
	mu.Lock()
	defer mu.Unlock()

	var alive []*Backend
	for _, b := range backends {
		if b.Alive {
			alive = append(alive, b)
		}
	}

	if len(alive) == 0 {
		return nil
	}

	b := alive[counter%len(alive)]
	counter++
	b.Connections++
	return b
}


func handler(w http.ResponseWriter, r *http.Request) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = r.RemoteAddr
	}

	mu.Lock()
	rateLimiter[ip]++
	count := rateLimiter[ip]
	mu.Unlock()

	if count > limitePerMinute {
		http.Error(w, "Too many requests", http.StatusTooManyRequests)
		return
	}

	backend := getNextBackend()
	if backend == nil {
		http.Error(w, "No healthy backends", http.StatusServiceUnavailable)
		return
	}

	defer func() {
		mu.Lock()
		backend.Connections--
		mu.Unlock()
	}()

	proxy := httputil.NewSingleHostReverseProxy(backend.URL)
	w.Write([]byte("Forwarded to " + backend.URL.String() + "\n"))
	proxy.ServeHTTP(w, r)
}

func startAdminAPI() {
	mux := http.NewServeMux()

	mux.HandleFunc("/admin/backends", listBackends)
	mux.HandleFunc("/admin/backends/add", addBackend)
	mux.HandleFunc("/admin/backends/remove", removeBackend)

	mux.HandleFunc("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "openapi.yaml")
	})

	mux.Handle(
		"/swagger/",
		httpSwagger.Handler(
			httpSwagger.URL("http://localhost:9091/openapi.yaml"),
		),
	)

	log.Println("Admin API listening on :9091")
	log.Fatal(http.ListenAndServe(":9091", mux))
}


func listBackends(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	json.NewEncoder(w).Encode(backends)
}

func addBackend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST required", http.StatusMethodNotAllowed)
		return
	}
	
	raw := r.URL.Query().Get("url")
	if raw == "" {
		http.Error(w, "missing url", http.StatusBadRequest)
		return
	}

	u, err := url.Parse(raw)
	if err != nil {
		http.Error(w, "invalid url", http.StatusBadRequest)
		return
	}

	mu.Lock()
	backends = append(backends, &Backend{URL: u, Alive: true})
	mu.Unlock()

	w.Write([]byte("backend added\n"))
}

func removeBackend(w http.ResponseWriter, r *http.Request) {
	raw := r.URL.Query().Get("url")
	if raw == "" {
		http.Error(w, "missing url", http.StatusBadRequest)
		return
	}
	mu.Lock()
	defer mu.Unlock()

	for i, b := range backends {
		if b.URL.String() == raw {
			backends = append(backends[:i], backends[i+1:]...)
			w.Write([]byte("backend removed\n"))
			return
		}
	}
	http.Error(w, "backend not found", http.StatusNotFound)
}


func main() {
	go resetRateLimiter()
	go startAdminAPI()
	http.HandleFunc("/", handler)

	log.Println("Proxy listening on :9090")
	log.Fatal(http.ListenAndServe(":9090", nil))
}

