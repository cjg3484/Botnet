package main

import (
	"bytes"
	"fmt"
	"github.com/kbinani/screenshot"
	"golang.org/x/sys/windows/registry"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

var opersys = runtime.GOOS

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
	} else if opersys == "linux" {
		data, err := ioutil.ReadFile("/etc/machine-id")
		if err != nil {
			panic(err)
		}
		id = string(data)
	} else {
		fmt.Println("unsupported OS")
		os.Exit(1)
	}
	return id
}

var idnum = getmachineid()

func call(urlPath, method string, filename string) error {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// New multipart writer.
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, err := writer.CreateFormFile("photo", filename)
	if err != nil {
		return err
	}
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	_, err = io.Copy(fw, file)
	if err != nil {
		return err
	}
	writer.Close()
	req, err := http.NewRequest(method, urlPath, bytes.NewReader(body.Bytes()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rsp, _ := client.Do(req)
	if rsp.StatusCode != http.StatusOK {
		log.Printf("Request failed with response code: %d", rsp.StatusCode)
	}
	return nil
}

func screengrab() {
	n := screenshot.NumActiveDisplays()
	var i int

	for i = 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)

		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			panic(err)
		}
		fileName := fmt.Sprintf("%s_%d.png", idnum, i)
		file, _ := os.Create(fileName)
		defer file.Close()
		png.Encode(file, img)

		//fmt.Printf("#%d : %v \"%s\"\n", i, bounds, fileName)
		switch opersys {
		case "windows":
			err = call("http://localhost:8081/upload", "POST", fileName)
		case "linux":
			err = call("http://127.0.1.1:8081/upload", "POST", fileName)
		case "default":
			fmt.Println("OS unsupported")
			os.Exit(1)
		}
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	screengrab()
}
