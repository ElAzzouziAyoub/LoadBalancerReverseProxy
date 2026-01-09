package main

import (
	"fmt"
	"net/http"
)

func handle(w http.ResponseWriter, r *http.Request) {

}

func main() {
	fmt.Print("Starting Server 2 : ")
	http.HandleFunc("/backend2",handle)
	http.ListenAndServe(":8082",nil)
}
