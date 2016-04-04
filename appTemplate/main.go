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

func main() {
	MakeActionRoutes()
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("client"))))
	//
	fmt.Println("Goster is running on http://127.0.0.1:8000")
	http.ListenAndServe(":8000", nil)
}
