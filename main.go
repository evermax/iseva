package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
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

	handler, err := NewJSONHandler(dbFile, staticGen)
	if err != nil {
		fmt.Printf("Problem when starting the server: %v\n", err)
		os.Exit(1)
	}
	http.HandleFunc("/", handler.ServeHTTP)
	fmt.Print("Starting server\n")
	http.ListenAndServe(":3000", nil)
}
