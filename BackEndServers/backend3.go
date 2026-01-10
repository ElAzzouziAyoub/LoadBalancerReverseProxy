package main

import (
	"fmt"
	"net/http"
)

func handle(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is backend 3 server"))
	switch r.Method {
	case http.MethodGet :
		fmt.Println("THis is a GET request")
	case http.MethodPost :
		fmt.Println("This is a Post request")
	}

}

func main() {
	fmt.Println("Starting Server 3 : ")
	http.HandleFunc("/",handle)
	http.ListenAndServe(":8083",nil)
}
