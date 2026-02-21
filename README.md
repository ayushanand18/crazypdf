# crazypdf

A generic, extensible PDF processing library for Go. Designed with a modular architecture where each capability (text extraction, structure analysis, etc.) lives in its own feature package under a shared core document model.

## Features

- **Text Extraction** — Extract text from PDFs with three layout modes:
  - **Simple** — Plain text, words joined by spaces, rows by newlines
  - **Raw** — Content stream order, preserving the internal PDF text order
  - **Physical** — Spatial layout preservation using x,y coordinates
- **Per-Page Access** — Access individual pages by index
- **CLI Tool** — Command-line utility with subcommand architecture
- **Pure Go** — No CGo dependencies, built on [ledongthuc/pdf](https://github.com/ledongthuc/pdf)

## Installation

### As a library

```bash
go get github.com/ayushanand18/crazypdf
```

### CLI tool

```bash
go install github.com/ayushanand18/crazypdf/cmd/crazypdf@latest
```

## Quick Start

### Library Usage

```go
package main

import (
    "fmt"
    "log"

    "github.com/ayushanand18/crazypdf/pkg/crazypdf"
    "github.com/ayushanand18/crazypdf/pkg/extract"
)

func main() {
    // Open a PDF document
    doc, err := crazypdf.Open("document.pdf")
    if err != nil {
        log.Fatal(err)
    }
    defer doc.Close()

    // Extract all text (simple mode)
    text, err := extract.Text(doc)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(text)
}
```

### Open from Bytes

```go
// Read PDF into memory
data, err := os.ReadFile("document.pdf")
if err != nil {
    log.Fatal(err)
}

// Open from byte slice
doc, err := crazypdf.OpenBytes(data)
if err != nil {
    log.Fatal(err)
}
defer doc.Close()
```

### Per-Page Extraction

```go
// Iterate over pages
for _, page := range doc.Pages() {
    text, err := extract.PageText(page)
    if err != nil {
        log.Printf("Page %d error: %v", page.Number, err)
        continue
    }
    fmt.Printf("--- Page %d ---\n%s\n", page.Number, text)
}
```

### Layout Modes

```go
// Physical layout — preserves spatial positioning
text, _ := extract.Text(doc, extract.WithLayout(extract.LayoutPhysical))

// Raw — content stream order
text, _ = extract.Text(doc, extract.WithLayout(extract.LayoutRaw))

// Custom page separator
text, _ = extract.Text(doc, extract.WithPageSeparator("\n---\n"))
```

### Encrypted PDFs

```go
doc, err := crazypdf.Open("encrypted.pdf", crazypdf.WithPassword("secret"))
```

## CLI Usage

```bash
# Extract text from a PDF (prints to stdout)
crazypdf text document.pdf

# Save to file
crazypdf text document.pdf output.txt

# Preserve physical layout
crazypdf text -layout document.pdf

# Raw content stream order
crazypdf text -raw document.pdf

# Specific pages
crazypdf text -pages 1-3 document.pdf
crazypdf text -pages 1,3,5 document.pdf

# Encrypted PDF
crazypdf text -password secret encrypted.pdf
```

## Architecture

```
crazypdf/
├── pkg/
│   ├── crazypdf/            # Core library package
│   │   ├── doc.go           # Package documentation
│   │   ├── document.go      # Document struct, Open/Close
│   │   ├── page.go          # Page struct, text accessors
│   │   ├── options.go       # Config, functional options
│   │   └── errors.go        # Shared error types
│   │
│   └── extract/             # Feature: Text Extraction
│       ├── text.go          # Text, PageText, AllPages
│       └── options.go       # LayoutMode, extraction options
│
├── internal/pdf/            # Internal PDF reader wrapper
│   └── reader.go            # Wraps ledongthuc/pdf
│
├── cmd/crazypdf/            # CLI tool
│   └── main.go              # Subcommand-based CLI
│
└── testdata/                # Test fixtures
    └── sample.pdf
```

### Adding New Features

The library is designed for extension. To add a new feature:

1. Create a new directory under `pkg/` (e.g., `pkg/structurize/`, `pkg/metadata/`)
2. Accept `*crazypdf.Document` or `*crazypdf.Page` as input
3. Use public accessor methods (`PlainText()`, `TextByRow()`, `StyledTexts()`, `ContentStream()`)
4. Add a new subcommand to `cmd/crazypdf/main.go`

### Planned Features

- **structurize** — Convert PDF structure into machine-readable format (headings, paragraphs, lists)
- **metadata** — Extract document metadata (title, author, keywords, dates)
- **tables** — Detect and extract tabular data

## API Reference

### Core Package (`pkg/crazypdf`)

| Type/Function | Description |
|---|---|
| `Open(path, ...Option) (*Document, error)` | Open a PDF file |
| `WithPassword(string) Option` | Set password for encrypted PDFs |
| `Document.NumPages() int` | Get page count |
| `Document.Page(index int) (*Page, error)` | Get page by 0-based index |
| `Document.Pages() []*Page` | Get all pages |
| `Document.Close() error` | Release resources |
| `Page.PlainText() (string, error)` | Get plain text from page |
| `Page.TextByRow() ([]TextRow, error)` | Get text organized by rows |
| `Page.StyledTexts() ([]StyledText, error)` | Get styled text with positions |
| `Page.ContentStream() ([]byte, error)` | Get raw content stream |

### Extract Package (`pkg/extract`)

| Type/Function | Description |
|---|---|
| `Text(doc, ...Option) (string, error)` | Extract text from entire document |
| `PageText(page, ...Option) (string, error)` | Extract text from single page |
| `AllPages(doc, ...Option) ([]string, error)` | Extract text from all pages |
| `WithLayout(LayoutMode) Option` | Set layout mode |
| `WithPageSeparator(string) Option` | Set page separator |
| `LayoutSimple` | Plain text extraction |
| `LayoutRaw` | Content stream order |
| `LayoutPhysical` | Spatial layout preservation |

## License

See [LICENSE](LICENSE) for details.
