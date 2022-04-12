package main

import (
	"bytes"
	"encoding/json"
	"golang.org/x/sys/windows/registry"
	"log"
	"net/http"
	"strings"
	"time"
)

type bot struct {
	Id     string `json:"bot_id"`
	Status string `json:"status"`
}

func getmachineid() string {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\SQMClient`, registry.QUERY_VALUE)
	if err != nil {
		log.Fatal(err)
	}
	defer k.Close()

	s, _, err := k.GetStringValue("MachineId")
	s2 := strings.Trim(s, "{}")
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(s2)

	return s2
}

func register() {
	idnum := getmachineid()

	b := bot{
		Id:     idnum,
		Status: "alive",
	}

	postBody, _ := json.Marshal(b)

	resp, err := http.Post("http://127.0.0.1:8081/register", "application/json", bytes.NewBuffer(postBody))
	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}

	//fmt.Println(resp)

	defer resp.Body.Close()
}

func heartbeat() {
	idnum := getmachineid()

	b := bot{
		Id:     idnum,
		Status: "alive",
	}

	postBody, _ := json.Marshal(b)

	resp, err := http.Post("http://127.0.0.1:8081/heartbeat", "application/json", bytes.NewBuffer(postBody))
	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}

	//fmt.Println(resp)

	defer resp.Body.Close()
}

func main() {

	register()

	for true {
		time.Sleep(3 * time.Second)
		heartbeat()
	}
}
