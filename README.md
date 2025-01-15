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

## Setup for docker

Go build

```bash
env GOOS=linux GOARCH=amd64 go build -o app .
```

Docker build

```bash
docker build -t app
```

Docker run

```bash
docker run --cpus=2 -p 8080:8080 app
```

## Setup for local

1. Install Go Dependencies:

```bash
go mod tidy
```

2. Install Chromium:

   - Ubuntu/Debian: `sudo apt install -y chromium`
   - Alpine: `apk add --no-cache chromium`

3. Run the Server:

```bash
go run main.go
```

## API Endpoint

**POST /generate-pdf**

- <b>Description</b>: Generates a PDF from an HTML template, merges it with two additional PDFs, and directly returns the result as a downloadable file (in memory).

- Example:

```bash
curl -X POST http://localhost:8080/generate/pdf-in-memory --output result.pdf
```

**POST /generate-pdf-save**

- <b>Description</b>: Generates a PDF from an HTML template, merges it with two additional PDFs, saves the merged result to disk, and then serves it to the client.

- Example:

```bash
curl -X POST http://localhost:8080/generate/pdf-to-disk --output result.pdf
```

## Note

This is a PoC for demonstration purposes and not intended for production use.
