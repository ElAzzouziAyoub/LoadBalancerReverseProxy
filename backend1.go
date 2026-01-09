package main

import (
	"fmt"
	"net/http"
)

func handle(w http.ResponseWriter, r *http.Request) {

}

func main() {
	fmt.Print("Starting Server 1 : ")
	http.HandleFunc("/backend1",handle)
	http.ListenAndServe(":8081",nil)
}
