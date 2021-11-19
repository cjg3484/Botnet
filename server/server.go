package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
)

// temporary directory location
// var srcDir = filepath.FromSlash("C:\\Users\\laugh\\Documents\\Github\\Botnet\\src")
var srcDir = filepath.FromSlash("H:\\GitHub\\Botnet\\src")

func main() {

	fmt.Printf("Starting server at port 8081\n")

	// return a `.pdf` file for `/pdf` route
	http.HandleFunc("/pdf", func(res http.ResponseWriter, req *http.Request) {
		http.ServeFile(res, req, filepath.Join(srcDir, "/files/PWNED.pdf"))
	})

	// start HTTP server with `http.DefaultServeMux` handler
	log.Fatal(http.ListenAndServe("localhost:8081", nil))
}
