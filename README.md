
# PDF Generator and Merger API (PoC)

This is a Proof of Concept (PoC) API built with Go to demonstrate generating PDFs from HTML templates and merging them with existing PDF files using chromedp and pdfcpu.

## Features

- Generate PDF dynamically from HTML templates.
- Merge generated PDFs with existing PDF files.
- In-memory processing for efficiency.

## Requirements

- Go (1.18+)
- Chromium or Google Chrome
- Libraries:
  - chromedp for HTML-to-PDF conversion.
  - pdfcpu for PDF merging.

## Setup

Go build

``` bash
env GOOS=linux GOARCH=amd64 go build -o app .
```

Docker build
``` bash
docker build -t app
```

Docker run
``` bash
docker run --cpus=2 -p 8080:8080 app
```

## API Endpoint

**POST /generate-pdf**

- Description: Generates a PDF from a template, merges it with two additional PDFs, and returns the result.

- Example:

```bash
curl -X POST http://localhost:8080/generate-pdf --output result.pdf
```

## Note

This is a PoC for demonstration purposes and not intended for production use.

