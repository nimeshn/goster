package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

var (
	AppDB *sql.DB
)

func Check(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func GetDB() *sql.DB {
	return AppDB
}

func main() {
	var err error
	AppDB, err = sql.Open("mysql",
		"root:shsemin123@tcp(127.0.0.1:3306)/sampleapp?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer AppDB.Close()
	err = AppDB.Ping()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("The database is accessible..")
	}
	MakeActionRoutes()
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("client"))))
	//
	fmt.Println("Goster is running on http://127.0.0.1:8000")
	http.ListenAndServe(":8000", nil)
}
