package main

import (
	"fmt"
	"log"
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
		log.Fatal(err)
	}
	return u
}

func handler(w http.ResponseWriter, r *http.Request) {
	//Figured there would a datarace here since http requests are sent in parallel
	mu.Lock()
	backend := backends[counter%len(backends)] //This line distributes http requests over all backend servers evenly
	counter++
	mu.Unlock()

	proxy := httputil.NewSingleHostReverseProxy(backend)
	w.Write([]byte("Your request was forwarded to Server "+ fmt.Sprintf("%d",counter%len(backends)+ 1) + "\n" ))
	proxy.ServeHTTP(w, r)
	fmt.Println(r.Header)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":9090", nil)
}

