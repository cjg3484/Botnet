package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	idnum := r1.Intn(10000)
	fmt.Print(idnum)

	postBody, _ := json.Marshal(map[string]string{
		"id":     string(idnum),
		"status": "alive",
	})
	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post("http://127.0.0.1:8081/register", "application/json", responseBody)
	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	log.Printf(sb)
}
