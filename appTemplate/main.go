package main

import (
	"fmt"
	"net/http"
)

func main() {
	MakeActionRoutes()
	//
	fmt.Println("Goster is running on http://127.0.0.1:8000")
	//
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("client"))))
	//
	http.ListenAndServe(":8000", nil)
}
