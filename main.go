package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/exec"

	"github.com/labstack/echo/v4"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func main() {
	e := echo.New()

	// Define the API endpoint
	e.POST("/generate-pdf", func(c echo.Context) error {
		// Step 1: Prepare the data for the template
		data := map[string]string{
			"someKey": "Hello, PDF!",
		}

		// Step 2: Process the HTML template
		htmlContent, err := processTemplate("templates/template.html", data)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		// Step 3: Convert HTML to PDF (using wkhtmltopdf)
		tempHTMLPath := "output/temp.html"
		tempPDFPath := "output/temp.pdf"
		if err := os.WriteFile(tempHTMLPath, []byte(htmlContent), 0644); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to write temp HTML file"})
		}

		if err := convertHTMLToPDF(tempHTMLPath, tempPDFPath); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to convert HTML to PDF"})
		}

		// Step 4: Merge with another PDF
		additionalPDFPath := "assets/additional.pdf"
		mergedPDFPath := "output/merged.pdf"
		if err := mergePDFs([]string{tempPDFPath, additionalPDFPath}, mergedPDFPath); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to merge PDFs"})
		}

		// Step 5: Read the merged PDF and return as byte stream
		mergedPDF, err := os.ReadFile(mergedPDFPath)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to read merged PDF"})
		}

		return c.Blob(http.StatusOK, "application/pdf", mergedPDF)
	})

	// Start the server
	e.Logger.Fatal(e.Start(":8080"))
}

// processTemplate replaces placeholders in the HTML template
func processTemplate(templatePath string, data map[string]string) (string, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// convertHTMLToPDF converts an HTML file to a PDF using wkhtmltopdf
func convertHTMLToPDF(htmlPath, pdfPath string) error {
	cmd := exec.Command("wkhtmltopdf", htmlPath, pdfPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running wkhtmltopdf: %v", err)
	}
	return nil
}

// mergePDFs merges multiple PDFs into one using pdfcpu
func mergePDFs(pdfPaths []string, outputPath string) error {
	// Configuration for pdfcpu (use nil for default configuration)
	conf := model.NewDefaultConfiguration()

	// Merge PDFs into the output file
	if err := api.MergeCreateFile(pdfPaths, outputPath, false, conf); err != nil {
		return fmt.Errorf("error merging PDFs: %v", err)
	}

	return nil
}
