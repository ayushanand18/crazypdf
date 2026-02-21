package crazypdf

import (
	internalpdf "github.com/ayushanand18/crazypdf/internal/pdf"
)

// Page represents a single page within a PDF document.
type Page struct {
	// Number is the 1-based page number.
	Number int

	doc *Document
}

// PlainText extracts plain text from this page with words joined by spaces
// and rows separated by newlines.
func (p *Page) PlainText() (string, error) {
	if p.doc.closed {
		return "", ErrDocumentClosed
	}
	return p.doc.reader.PagePlainText(p.Number)
}

// TextByRow returns text organized by rows with position information.
func (p *Page) TextByRow() ([]internalpdf.TextRow, error) {
	if p.doc.closed {
		return nil, ErrDocumentClosed
	}
	return p.doc.reader.PageTextByRow(p.Number)
}

// StyledTexts returns text elements with font and position information.
func (p *Page) StyledTexts() ([]internalpdf.StyledText, error) {
	if p.doc.closed {
		return nil, ErrDocumentClosed
	}
	return p.doc.reader.PageStyledTexts(p.Number)
}

// ContentStream returns the raw PDF content stream bytes for this page.
func (p *Page) ContentStream() ([]byte, error) {
	if p.doc.closed {
		return nil, ErrDocumentClosed
	}
	return p.doc.reader.PageContentStream(p.Number)
}

// PhysicalLayoutText extracts text preserving spatial positioning on the page.
// pageWidth is the page width in PDF points (default 612 for US Letter).
func (p *Page) PhysicalLayoutText(pageWidth float64) (string, error) {
	if p.doc.closed {
		return "", ErrDocumentClosed
	}
	return p.doc.reader.PhysicalLayoutText(p.Number, pageWidth)
}
