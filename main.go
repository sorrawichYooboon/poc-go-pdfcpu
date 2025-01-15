package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
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

		// Step 3: Convert HTML to PDF using chromedp
		// tempPDFPath := "output/temp.pdf"
		pdfBuffer, err := convertHTMLToPDFWithChromedp(htmlContent)
		if err != nil {
			fmt.Println(err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to convert HTML to PDF"})
		}

		// Step 4: Merge with another PDF
		additionalPDFPath := "assets/test1.pdf"
		additional2PDFPath := "assets/test2.pdf"
		// mergedPDFPath := "output/merged.pdf"

		// Read additional PDF files into memory
		additionalPDF, err := os.ReadFile(additionalPDFPath)
		if err != nil {
			return fmt.Errorf("failed to read additional PDF file: %v", err)
		}

		additional2PDF, err := os.ReadFile(additional2PDFPath)
		if err != nil {
			return fmt.Errorf("failed to read second additional PDF file: %v", err)
		}

		b, err := mergePDFs(pdfBuffer, additionalPDF, additional2PDF)
		if err != nil {
			return fmt.Errorf("error merging PDFs: %v", err)
		}


		return c.Blob(http.StatusOK, "application/pdf", b.Bytes())
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

// convertHTMLToPDFWithChromedp converts HTML content to a PDF using chromedp
func convertHTMLToPDFWithChromedp(htmlContent string) ([]byte, error) {
	// Create a new Chrome context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Buffer to store the generated PDF
	var pdfBuffer []byte

	// Run chromedp tasks
	err := chromedp.Run(ctx,
		// Navigate to the HTML content
		chromedp.Navigate(`data:text/html,`+htmlContent),

		// Generate PDF from the HTML content
		chromedp.EmulateViewport(1280, 1024), // Optional: Adjust viewport for rendering
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdfBuffer, _, err = page.PrintToPDF().WithPrintBackground(true).Do(ctx)
			return err
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("error generating PDF with chromedp: %v", err)
	}

	return pdfBuffer, nil
}

// mergePDFs merges multiple PDFs into one using pdfcpu
func mergePDFs(b ...[]byte) (*bytes.Buffer, error) {
	// Configuration for pdfcpu (use nil for default configuration)
	conf := model.NewDefaultConfiguration()
	readers := []io.ReadSeeker{}
	// Create readers for the PDFs
	for _, v := range b {
		readers = append(readers, bytes.NewReader(v))
	}

	// Buffer to hold the merged PDF data
	tempw := new(bytes.Buffer)

	// Merge PDFs into the output file
	if err := api.MergeRaw(readers, tempw, false, conf); err != nil {
		return nil, fmt.Errorf("error merging PDFs: %v", err)
	}
	return tempw, nil
}
