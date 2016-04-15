package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
)

func rootHandler(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(rw, "Hello Goster!")

}

func main() {
	model := Model{
		Name:        "post",
		DisplayName: "User Posts",
		Fields: []*Field{
			&Field{
				Name:        "hdr",
				DisplayName: "Title",
				Type:        String,
				Validator: &FieldValidation{
					Required: true,
					MaxLen:   50,
				},
			},
			&Field{
				Name:        "descript",
				DisplayName: "Message",
				Type:        String,
				Validator: &FieldValidation{
					Required: true,
					MaxLen:   200,
				},
			},
			&Field{
				Name:        "location",
				DisplayName: "Location",
				Type:        String,
				Validator: &FieldValidation{
					Required: true,
					MaxLen:   255,
					MinLen:   5,
				},
			},
		},
	}

	workDir, err := os.Getwd()
	Check(err)
	app := CreateNewApp("SampleApp", "Sample Application", "Bitwinger", "V1.0",
		path.Join(path.Dir(workDir), "SampleApp"))
	app.AddModel(&model)
	app.MakeClient()

	app.MakeServer()

	json.Marshal(app)
	//vals, _ := json.Marshal(app)
	//fmt.Println(string(vals))

	/*

		fmt.Println("Unmarshalling now")
		str := `{"name": "Hi2", "DisplayName": "Name2"}`
		jsonModel := Model{}
		json.Unmarshal([]byte(str), &jsonModel)
		fmt.Println(jsonModel)
	*/
	fmt.Println("Goster is running on http://127.0.0.1:8000")

	//
	http.Handle("/client/", http.StripPrefix("/client/", http.FileServer(http.Dir("client"))))
	//
	//http.HandleFunc("/", rootHandler)
	//http.ListenAndServe(":8000", nil)
}
