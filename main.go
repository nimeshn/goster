package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func rootHandler(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(rw, "Hello Goster!")

}

func main() {

	model := Model{
		Name:        "Hi",
		DisplayName: "Name",
		Fields:      []*Field{&Field{Validator: &FieldValidation{}}}}
	vals, _ := json.Marshal(model)
	fmt.Println(string(vals))

	fmt.Println("Unmarshalling now")
	str := `{"name": "Hi2", "DisplayName": "Name2"}`
	jsonModel := Model{}
	json.Unmarshal([]byte(str), &jsonModel)
	fmt.Println(jsonModel)

	fmt.Println("Goster is running on http://127.0.0.1:8000")
	//
	http.Handle("/client/", http.StripPrefix("/client/", http.FileServer(http.Dir("client"))))
	//
	//http.HandleFunc("/", rootHandler)
	//http.ListenAndServe(":8000", nil)
}
