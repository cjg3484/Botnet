package main

import (
	"github.com/jung-kurt/gofpdf"
	"log"
	"os/exec"
	"path/filepath"
)

func main() {
	var fileName = filepath.FromSlash("C:\\Users\\laugh\\GolandProjects\\Botnet\\src\\files\\hello.pdf")

	//create pdf
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Hello, world")
	err := pdf.OutputFileAndClose(fileName)
	if err != nil {
		return
	}

	//open pdf
	cmd := exec.Command("cmd.exe", "/C", "start", fileName)
	if err := cmd.Run(); err != nil {
		log.Println("Error:", err)
	}

	//download client

}
