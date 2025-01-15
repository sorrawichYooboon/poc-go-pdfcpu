package poc

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/labstack/echo/v4"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func GeneratePDFToDisk(c echo.Context) error {
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
	tempPDFPath := "output/html_to_pdf.pdf"
	if err := convertHTMLToPDFWithChromedpSave(htmlContent, tempPDFPath); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to convert HTML to PDF"})
	}

	// Step 4: Merge with another PDF
	additionalPDFPath := "assets/test1.pdf"
	additional2PDFPath := "assets/test2.pdf"
	mergedPDFPath := "output/merged.pdf"
	if err := mergePDFsSave([]string{tempPDFPath, additionalPDFPath, additional2PDFPath}, mergedPDFPath); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to merge PDFs"})
	}

	// Step 5: Read the merged PDF and return as byte stream
	mergedPDF, err := os.ReadFile(mergedPDFPath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to read merged PDF"})
	}

	return c.Blob(http.StatusOK, "application/pdf", mergedPDF)
}

// convertHTMLToPDFWithChromedp converts HTML content to a PDF using chromedp
func convertHTMLToPDFWithChromedpSave(htmlContent, pdfPath string) error {
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
		return fmt.Errorf("error generating PDF with chromedp: %v", err)
	}

	// Write the PDF to the specified file
	err = os.WriteFile(pdfPath, pdfBuffer, 0644)
	if err != nil {
		return fmt.Errorf("error saving PDF: %v", err)
	}

	return nil
}

// mergePDFs merges multiple PDFs into one using pdfcpu
func mergePDFsSave(pdfPaths []string, outputPath string) error {
	// Configuration for pdfcpu (use nil for default configuration)
	conf := model.NewDefaultConfiguration()

	// Merge PDFs into the output file
	if err := api.MergeCreateFile(pdfPaths, outputPath, false, conf); err != nil {
		return fmt.Errorf("error merging PDFs: %v", err)
	}

	return nil
}
