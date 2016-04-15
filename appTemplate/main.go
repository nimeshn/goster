package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os/exec"
	"runtime"
)

var (
	AppDB      *sql.DB
	serverUrl  string = "http://127.0.0.1:%d"
	portNo     int    = 8000
	dbUser     string = "root"
	dbPassword string = "shsemin123"
	dbName     string = "sampleapp"
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

func OpenURL(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows", "darwin":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	serverUrl = fmt.Sprintf(serverUrl, portNo)
	//Connect to the Database
	var err error
	AppDB, err = sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?parseTime=true", dbUser, dbPassword, dbName))
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
	fmt.Printf("Goster is starting on %s...", serverUrl)
	OpenURL(serverUrl)
	http.ListenAndServe(fmt.Sprintf(":%d", portNo), nil)
}
