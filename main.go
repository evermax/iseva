package main

import (
	"flag"
	"fmt"
	"net/http"
)

const (
	dbPath = "db.json"
)

var dbFile string
var staticGen bool

func init() {
	flag.StringVar(&dbFile, "db", dbPath, "Specify the path of the file in which the JSON is. The default value is db.json")
	flag.BoolVar(&staticGen, "s", false, "Specify if you want the JSON file to be loaded on every request or imported in memory and statically serve. This means the random values will be set for the time the program runs. The default value is false")
}

func main() {
	flag.Parse()

	handler := JSONHandler{DB: dbFile, IsStatic: staticGen}
	handler.Init()
	http.HandleFunc("/", handler.ServeHTTP)
	fmt.Print("Starting server\n")
	http.ListenAndServe(":3000", nil)
}
