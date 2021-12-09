package main

import (
	"os/exec"

	"github.com/jung-kurt/gofpdf"
)

func main() {
	//var srcDir = filepath.FromSlash("H:\\GolandProjects\\Botnet\\src")

	//create pdf
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Hello, world")
	pdf.OutputFileAndClose("hello.pdf")

	//open pdf
	exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", "hello.pdf")

	//download client
}
