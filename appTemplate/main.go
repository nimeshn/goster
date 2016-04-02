package main

import (
	"fmt"
	"net/http"
)

func Check(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func rootHandler(rw http.ResponseWriter, req *http.Request) {
	fmt.Println(req)
	fmt.Fprintf(rw, "Hello Goster!")
}

func main() {
	MakeActionRoutes()
	//
	fmt.Println("Goster is running on http://127.0.0.1:8000")
	//
	//http.Handle("/api/", http.StripPrefix("/api/", rootHandler))
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("client"))))
	//
	http.ListenAndServe(":8000", nil)
}
