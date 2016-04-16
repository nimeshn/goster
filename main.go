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

	app := CreateNewApp("SampleApp", "Sample Application", "Bitwinger", "V1.0", 9000)
	fmt.Println(app.AppDir)
	app.AddModel(&model)
	SaveAppSettings(app)
	app.MakeClient()
	app.MakeServer()
	app.InstallAndRunApp()

	json.Marshal(app)
	fmt.Println("Goster is running on http://127.0.0.1:8000")
	http.Handle("/client/", http.StripPrefix("/client/", http.FileServer(http.Dir("client"))))
	//
	//http.HandleFunc("/", rootHandler)
	//http.ListenAndServe(":8000", nil)
}
