package main

import (
	"fmt"
	"net/http"
)


func HandleGetUsers(w http.ResponseWriter, r *http.Request) {


}

func HandlePostUsers(w http.ResponseWriter, r *http.Request) {

}

func handler(w http.ResponseWriter, req *http.Request) {

	fmt.Println("Connection established with ",req.RemoteAddr)
	w.Write([]byte("Connection established with server \n"))
	switch req.Method {
	case http.MethodGet:
		HandleGetUsers(w, req)
	case http.MethodPost:
		HandlePostUsers(w, req)
	default:
		w.WriteHeader(http.StatusNotImplemented)
	}

}

func main() {
	fmt.Println("Started Server 1 : ")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8081", nil)

}
