package crazypdf

import (
	"fmt"

	internalpdf "github.com/ayushanand18/crazypdf/internal/pdf"
)

// Document represents an opened PDF document.
// It is the central type that all feature modules operate on.
type Document struct {
	filePath string
	reader   *internalpdf.Reader
	pages    []*Page
	config   *Config
	closed   bool
}

// Open opens a PDF file from disk and returns a Document ready for processing.
func Open(filePath string, opts ...Option) (*Document, error) {
	cfg := applyOptions(opts)

	reader, err := internalpdf.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidPDF, err)
	}

	doc := &Document{
		filePath: filePath,
		reader:   reader,
		config:   cfg,
	}

	// Build page list
	numPages := reader.NumPages()
	doc.pages = make([]*Page, numPages)
	for i := 0; i < numPages; i++ {
		doc.pages[i] = &Page{
			Number: i + 1, // 1-based page number
			doc:    doc,
		}
	}

	return doc, nil
}

// OpenBytes opens a PDF from a byte slice and returns a Document ready for processing.
func OpenBytes(data []byte, opts ...Option) (*Document, error) {
	cfg := applyOptions(opts)

	reader, err := internalpdf.OpenBytes(data)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidPDF, err)
	}

	doc := &Document{
		filePath: "",
		reader:   reader,
		config:   cfg,
	}

	// Build page list
	numPages := reader.NumPages()
	doc.pages = make([]*Page, numPages)
	for i := 0; i < numPages; i++ {
		doc.pages[i] = &Page{
			Number: i + 1, // 1-based page number
			doc:    doc,
		}
	}

	return doc, nil
}

// NumPages returns the total number of pages in the document.
func (d *Document) NumPages() int {
	return len(d.pages)
}

// Page returns a specific page by 0-based index.
func (d *Document) Page(index int) (*Page, error) {
	if d.closed {
		return nil, ErrDocumentClosed
	}
	if index < 0 || index >= len(d.pages) {
		return nil, fmt.Errorf("%w: requested %d, document has %d pages", ErrPageOutOfRange, index, len(d.pages))
	}
	return d.pages[index], nil
}

// Pages returns all pages in the document.
func (d *Document) Pages() []*Page {
	return d.pages
}

// FilePath returns the file path of the opened document.
func (d *Document) FilePath() string {
	return d.filePath
}

// Close releases all resources held by the document.
func (d *Document) Close() error {
	if d.closed {
		return nil
	}
	d.closed = true
	if d.reader != nil {
		return d.reader.Close()
	}
	return nil
}

// Reader returns the internal PDF reader for advanced operations.
// This is intended for use by feature modules (e.g., extract).
func (d *Document) Reader() *internalpdf.Reader {
	return d.reader
}

// IsClosed returns whether the document has been closed.
func (d *Document) IsClosed() bool {
	return d.closed
}
