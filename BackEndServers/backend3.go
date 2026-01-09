package main

import (
	"fmt"
	"net/http"
)

func handle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request received from",r.RemoteAddr)

}

func main() {
	fmt.Println("Starting Server 3 : ")
	http.HandleFunc("/backend3",handle)
	http.ListenAndServe(":8083",nil)
}
