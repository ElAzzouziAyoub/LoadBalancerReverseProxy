package main

import (
	"fmt"
	"net/http"
)

func handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet :
		fmt.Println("THis is a GET request")
	case http.MethodPost :
		fmt.Println("This is a Post request")
	}
}

func main() {
	fmt.Print("Starting Server 2 : ")
	http.HandleFunc("/",handle)
	http.ListenAndServe(":8082",nil)
}
