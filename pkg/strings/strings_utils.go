package strings

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

// extractRawStrings extracts raw string content from PDF bytes
func (e *pdfStringExtractor) extractRawStrings(pdfBytes []byte) ([]string, error) {
	var strings []string

	// Convert bytes to string for regex processing
	pdfContent := string(pdfBytes)

	// Regex patterns to match PDF string literals
	patterns := []*regexp.Regexp{
		// Match literal strings: (string content)
		regexp.MustCompile(`\((.*?)\)`),
		// Match hex strings: <0123456789ABCDEF>
		regexp.MustCompile(`<([0-9A-Fa-f\s]+)>`),
		// Match text in content streams: Tj, TJ operators
		regexp.MustCompile(`\((.*?)\)\s*T[Jj]`),
	}

	for _, pattern := range patterns {
		matches := pattern.FindAllStringSubmatch(pdfContent, -1)
		for _, match := range matches {
			if len(match) > 1 && match[1] != "" {
				strings = append(strings, match[1])
			}
		}
	}

	return strings, nil
}

// processStrings cleans and filters extracted strings
func (e *pdfStringExtractor) processStrings(rawStrings []string) []string {
	var processed []string
	seen := make(map[string]bool)

	for _, str := range rawStrings {
		// Clean the string
		cleaned := e.cleanString(str)

		// Skip empty strings and very short ones (likely artifacts)
		if cleaned == "" || len(cleaned) < 2 {
			continue
		}

		// Remove duplicates
		if !seen[cleaned] {
			seen[cleaned] = true
			processed = append(processed, cleaned)
		}
	}

	return processed
}

// cleanString performs various cleaning operations on extracted strings
func (e *pdfStringExtractor) cleanString(str string) string {
	// Remove escape sequences and PDF-specific formatting
	cleaned := strings.ReplaceAll(str, "\\(", "(")
	cleaned = strings.ReplaceAll(cleaned, "\\)", ")")
	cleaned = strings.ReplaceAll(cleaned, "\\n", "\n")
	cleaned = strings.ReplaceAll(cleaned, "\\r", "\r")
	cleaned = strings.ReplaceAll(cleaned, "\\t", "\t")

	// Remove hex string artifacts
	cleaned = strings.TrimSpace(cleaned)

	// Remove common PDF control characters
	cleaned = regexp.MustCompile(`[\x00-\x1F\x7F]`).ReplaceAllString(cleaned, "")

	return cleaned
}

// extractStringsFromPage extracts strings from a specific page using pdfcpu's capabilities
func (e *pdfStringExtractor) extractStringsFromPage(ctx *model.Context, pageNumber int) ([]string, error) {
	var strings []string

	// Get page content
	content, err := api.ExtractPage(ctx, pageNumber)
	if err != nil {
		return nil, err
	}

	// read in chunks of 1024 bytes
	contentStr := bytes.Buffer{}
	b := make([]byte, 1024)
	for {
		n, err := content.Read(b)
		if err != nil || n < 1024 {
			break
		}
		contentStr.Write(b)
	}

	// Simple extraction from content streams
	// This is a basic approach - will parse the content stream later on
	textPattern := regexp.MustCompile(`\((.*?)\)\s*T[Jj]`)
	matches := textPattern.FindAllStringSubmatch(contentStr.String(), -1)

	for _, match := range matches {
		if len(match) > 1 {
			strings = append(strings, match[1])
		}
	}

	return strings, nil
}
