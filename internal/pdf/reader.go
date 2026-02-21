// Package pdf provides an internal wrapper around the ledongthuc/pdf library
// for reading PDF files and accessing page content.
package pdf

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"

	gopdf "github.com/ledongthuc/pdf"
)

// Reader wraps the ledongthuc/pdf reader and manages the underlying file handle.
type Reader struct {
	file   *os.File
	reader *gopdf.Reader
}

// OpenFile opens a PDF file from disk and returns a Reader.
func OpenFile(filePath string) (*Reader, error) {
	f, r, err := gopdf.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open PDF: %w", err)
	}
	return &Reader{file: f, reader: r}, nil
}

// OpenBytes opens a PDF from a byte slice and returns a Reader.
func OpenBytes(data []byte) (*Reader, error) {
	r, err := gopdf.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("failed to open PDF from bytes: %w", err)
	}
	return &Reader{file: nil, reader: r}, nil
}

// Close closes the underlying file handle.
func (r *Reader) Close() error {
	if r.file != nil {
		return r.file.Close()
	}
	return nil
}

// NumPages returns the total number of pages in the PDF.
func (r *Reader) NumPages() int {
	return r.reader.NumPage()
}

// PlainText extracts all plain text from the entire document.
func (r *Reader) PlainText() (string, error) {
	textReader, err := r.reader.GetPlainText()
	if err != nil {
		return "", fmt.Errorf("failed to get plain text: %w", err)
	}
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(textReader); err != nil {
		return "", fmt.Errorf("failed to read text: %w", err)
	}
	return buf.String(), nil
}

// PagePlainText extracts plain text from a specific page (1-based index).
// It uses X-position and font size data to intelligently merge adjacent
// glyph groups that belong to the same word, only inserting spaces where
// there is a genuine gap between words.
func (r *Reader) PagePlainText(pageNum int) (string, error) {
	page := r.reader.Page(pageNum)
	if page.V.IsNull() {
		return "", fmt.Errorf("page %d is null", pageNum)
	}

	rows, err := page.GetTextByRow()
	if err != nil {
		return "", fmt.Errorf("failed to get text for page %d: %w", pageNum, err)
	}

	var buf bytes.Buffer
	for i, row := range rows {
		if len(row.Content) == 0 {
			if i < len(rows)-1 {
				buf.WriteString("\n")
			}
			continue
		}

		// Sort content items by X position within this row
		items := make([]gopdf.Text, len(row.Content))
		copy(items, row.Content)
		sort.Slice(items, func(a, b int) bool {
			return items[a].X < items[b].X
		})

		// Adaptive character width estimation.
		// Compute per-character advance (gap / len(prev.S)) for each
		// consecutive pair. The minimum gives us a reliable estimate
		// of the actual character width (from contiguous glyph pairs).
		minCharWidth := 0.0
		for j := 1; j < len(items); j++ {
			prev := items[j-1]
			curr := items[j]
			if len(prev.S) > 0 {
				advancePerChar := (curr.X - prev.X) / float64(len(prev.S))
				if advancePerChar > 0 && (minCharWidth == 0 || advancePerChar < minCharWidth) {
					minCharWidth = advancePerChar
				}
			}
		}

		// Fallback: use fontSize * 0.6 if we couldn't compute from positions
		if minCharWidth == 0 {
			fontSize := items[0].FontSize
			if fontSize <= 0 {
				fontSize = 12
			}
			minCharWidth = fontSize * 0.6
		}

		buf.WriteString(items[0].S)

		for j := 1; j < len(items); j++ {
			prev := items[j-1]
			curr := items[j]

			// Expected end position of previous item if characters were contiguous
			prevEndX := prev.X + float64(len(prev.S))*minCharWidth
			gap := curr.X - prevEndX

			// If the gap exceeds half a character width, it's a word space.
			// This threshold accounts for kerning variations while still
			// catching genuine word separations.
			if gap > minCharWidth*0.5 {
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

// TextRow represents a row of text with its vertical position.
type TextRow struct {
	Position int64
	Words    []TextWord
}

// TextWord represents a word of text with its horizontal position and styling.
type TextWord struct {
	S        string
	X        float64
	Y        float64
	Font     string
	FontSize float64
}

// PageTextByRow returns text organized by rows for a specific page (1-based index).
func (r *Reader) PageTextByRow(pageNum int) ([]TextRow, error) {
	page := r.reader.Page(pageNum)
	if page.V.IsNull() {
		return nil, fmt.Errorf("page %d is null", pageNum)
	}

	rows, err := page.GetTextByRow()
	if err != nil {
		return nil, fmt.Errorf("failed to get text rows for page %d: %w", pageNum, err)
	}

	var result []TextRow
	for _, row := range rows {
		tr := TextRow{Position: row.Position}
		for _, word := range row.Content {
			tr.Words = append(tr.Words, TextWord{
				S:        word.S,
				X:        word.X,
				Y:        word.Y,
				Font:     word.Font,
				FontSize: word.FontSize,
			})
		}
		result = append(result, tr)
	}
	return result, nil
}

// StyledText represents text with font and position information.
type StyledText struct {
	Text     string
	X        float64
	Y        float64
	Font     string
	FontSize float64
}

// PageStyledTexts returns styled text elements for a specific page (1-based index).
// The returned texts include position and font information.
func (r *Reader) PageStyledTexts(pageNum int) ([]StyledText, error) {
	page := r.reader.Page(pageNum)
	if page.V.IsNull() {
		return nil, fmt.Errorf("page %d is null", pageNum)
	}

	rows, err := page.GetTextByRow()
	if err != nil {
		return nil, fmt.Errorf("failed to get styled texts for page %d: %w", pageNum, err)
	}

	var result []StyledText
	for _, row := range rows {
		for _, word := range row.Content {
			result = append(result, StyledText{
				Text:     word.S,
				X:        word.X,
				Y:        word.Y,
				Font:     word.Font,
				FontSize: word.FontSize,
			})
		}
	}
	return result, nil
}

// PageContentStream returns the raw content stream bytes for a page (1-based).
func (r *Reader) PageContentStream(pageNum int) ([]byte, error) {
	page := r.reader.Page(pageNum)
	if page.V.IsNull() {
		return nil, fmt.Errorf("page %d is null", pageNum)
	}

	content := page.V.Key("Contents")
	if content.Kind() == gopdf.Null {
		return nil, nil
	}

	reader := content.Reader()
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, reader); err != nil {
		return nil, fmt.Errorf("failed to read content stream for page %d: %w", pageNum, err)
	}
	return buf.Bytes(), nil
}

// PhysicalLayoutText extracts text preserving physical positioning for a page (1-based).
// It uses x,y coordinates to reconstruct the spatial layout of text on the page.
func (r *Reader) PhysicalLayoutText(pageNum int, pageWidth float64) (string, error) {
	styledTexts, err := r.PageStyledTexts(pageNum)
	if err != nil {
		return "", err
	}

	if len(styledTexts) == 0 {
		return "", nil
	}

	// Sort by Y (descending — PDF origin is bottom-left), then by X
	sort.Slice(styledTexts, func(i, j int) bool {
		if styledTexts[i].Y != styledTexts[j].Y {
			return styledTexts[i].Y > styledTexts[j].Y
		}
		return styledTexts[i].X < styledTexts[j].X
	})

	// Group texts by approximate Y position (same line if within tolerance)
	const yTolerance = 2.0
	type line struct {
		y     float64
		texts []StyledText
	}

	var lines []line
	var currentLine *line

	for _, st := range styledTexts {
		if currentLine == nil || abs(currentLine.y-st.Y) > yTolerance {
			lines = append(lines, line{y: st.Y})
			currentLine = &lines[len(lines)-1]
		}
		currentLine.texts = append(currentLine.texts, st)
	}

	// Determine column width: use average character width or default
	if pageWidth <= 0 {
		pageWidth = 612 // default US Letter width in points
	}
	charsPerLine := 80
	charWidth := pageWidth / float64(charsPerLine)

	var buf bytes.Buffer
	for i, ln := range lines {
		// Sort texts in this line by X position
		sort.Slice(ln.texts, func(a, b int) bool {
			return ln.texts[a].X < ln.texts[b].X
		})

		// Build the line with spacing
		lineChars := make([]byte, charsPerLine)
		for i := range lineChars {
			lineChars[i] = ' '
		}

		for _, st := range ln.texts {
			col := int(st.X / charWidth)
			if col < 0 {
				col = 0
			}
			if col >= charsPerLine {
				col = charsPerLine - 1
			}
			for ci, ch := range []byte(st.Text) {
				pos := col + ci
				if pos < charsPerLine {
					lineChars[pos] = ch
				}
			}
		}

		// Trim trailing spaces and write
		lineStr := string(lineChars)
		trimmed := trimRight(lineStr)
		buf.WriteString(trimmed)
		if i < len(lines)-1 {
			buf.WriteString("\n")
		}
	}

	return buf.String(), nil
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func trimRight(s string) string {
	i := len(s) - 1
	for i >= 0 && s[i] == ' ' {
		i--
	}
	return s[:i+1]
}
