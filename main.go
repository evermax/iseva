package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/evermax/iseva/env"
)

const (
	dbPath = "db.json"
)

var dbFile string

func init() {
	flag.StringVar(&dbFile, "db", dbPath, "Specify the path of the file in which the JSON is. The default value is db.json")
}

func main() {
	flag.Parse()

	env := env.Interface{DB: dbFile}
	http.HandleFunc("/", env.Handler)
	fmt.Print("Starting server\n")
	http.ListenAndServe(":3000", nil)
}
