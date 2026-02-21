
// Package crazypdf is a generic PDF processing library for Go.
//
// It provides a modular architecture where the core package defines the
// Document and Page types, and feature modules (like extract) operate on them.
//
// # Quick Start
//
// Open a PDF and extract text:
//
//	doc, err := crazypdf.Open("document.pdf")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer doc.Close()
//
//	// Use the extract package for text extraction
//	text, err := extract.Text(doc)
//	fmt.Println(text)
//
// # Architecture
//
// The library is organized into:
//   - pkg/crazypdf: Core types — Document, Page, Config, errors
//   - pkg/extract: Text extraction with multiple layout modes
//   - cmd/crazypdf: CLI tool with subcommand architecture
//
// Feature modules accept *Document or *Page and use their public methods
// (PlainText, TextByRow, StyledTexts, ContentStream, PhysicalLayoutText)
// to access page content without reaching into private fields.
//
// # Planned Features
//
//   - structurize: Convert PDF structure into machine-readable format
//   - metadata: Extract document metadata (title, author, keywords, etc.)
//   - tables: Detect and extract tabular data
package crazypdf
