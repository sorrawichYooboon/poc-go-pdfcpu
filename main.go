package main

import (
	"github.com/labstack/echo/v4"
	"github.com/sorrawichYooboon/poc-go-pdfcpu/poc"
)

func main() {
	e := echo.New()

	// In Memory PDF Generation
	e.POST("/generate/pdf-in-memory", poc.GeneratePDFInMemory)
	// Save PDF to Disk
	e.POST("/generate/pdf-to-disk", poc.GeneratePDFToDisk)

	// Start the server
	e.Logger.Fatal(e.Start(":8080"))
}
