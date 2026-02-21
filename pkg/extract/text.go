// Package extract provides text extraction capabilities for PDF documents.
//
// It operates on crazypdf.Document and crazypdf.Page types, supporting
// multiple layout modes for different text extraction needs.
package extract

import (
	"fmt"
	"sort"
	"strings"

	internalpdf "github.com/ayushanand18/crazypdf/internal/pdf"
	"github.com/ayushanand18/crazypdf/pkg/crazypdf"
)

// Text extracts text from the entire document, joining pages with the
// configured page separator. Default layout is LayoutSimple.
func Text(doc *crazypdf.Document, opts ...Option) (string, error) {
	if doc.IsClosed() {
		return "", crazypdf.ErrDocumentClosed
	}

	pages, err := AllPages(doc, opts...)
	if err != nil {
		return "", err
	}

	cfg := applyOptions(opts)
	return strings.Join(pages, cfg.PageSeparator), nil
}

// PageText extracts text from a single page.
func PageText(page *crazypdf.Page, opts ...Option) (string, error) {
	cfg := applyOptions(opts)

	switch cfg.Layout {
	case LayoutSimple:
		return page.PlainText()
	case LayoutRaw:
		return extractRawText(page)
	case LayoutPhysical:
		return page.PhysicalLayoutText(cfg.PageWidth)
	default:
		return page.PlainText()
	}
}

// AllPages extracts text from all pages, returning a slice with one entry per page.
func AllPages(doc *crazypdf.Document, opts ...Option) ([]string, error) {
	if doc.IsClosed() {
		return nil, crazypdf.ErrDocumentClosed
	}

	pages := doc.Pages()
	result := make([]string, 0, len(pages))

	for _, page := range pages {
		text, err := PageText(page, opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to extract text from page %d: %w", page.Number, err)
		}
		result = append(result, text)
	}

	return result, nil
}

// extractRawText extracts text in content stream order.
// This uses the row-based extraction from the reader which preserves
// the order text appears in the content stream. It uses X-position
// and font size data to intelligently merge adjacent glyph groups.
func extractRawText(page *crazypdf.Page) (string, error) {
	rows, err := page.TextByRow()
	if err != nil {
		return "", err
	}

	var buf strings.Builder
	for i, row := range rows {
		if len(row.Words) == 0 {
			if i < len(rows)-1 {
				buf.WriteString("\n")
			}
			continue
		}

		// Sort words by X position within the row
		words := make([]internalpdf.TextWord, len(row.Words))
		copy(words, row.Words)
		sort.Slice(words, func(a, b int) bool {
			return words[a].X < words[b].X
		})

		buf.WriteString(words[0].S)

		for j := 1; j < len(words); j++ {
			prev := words[j-1]
			curr := words[j]

			fontSize := prev.FontSize
			if fontSize <= 0 {
				fontSize = curr.FontSize
			}
			if fontSize <= 0 {
				fontSize = 12
			}
			avgCharWidth := fontSize * 0.5
			prevEndX := prev.X + float64(len(prev.S))*avgCharWidth

			gap := curr.X - prevEndX

			if gap > avgCharWidth*0.3 {
				buf.WriteString(" ")
			}
			buf.WriteString(curr.S)
		}

		if i < len(rows)-1 {
			buf.WriteString("\n")
		}
	}
	return buf.String(), nil
}
