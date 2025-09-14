package strings

import (
	"fmt"
	"log"
	"os"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

// ExtractAllStrings extracts all string literals from the PDF
func (e *pdfStringExtractor) ExtractAllStrings() ([]string, error) {
	// Read the PDF file as bytes
	pdfBytes, err := os.ReadFile(e.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read PDF file: %w", err)
	}

	// Extract raw content
	rawStrings, err := e.extractRawStrings(pdfBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to extract raw strings: %w", err)
	}

	// Process and filter strings
	processedStrings := e.processStrings(rawStrings)

	return processedStrings, nil
}

// ExtractStringsWithContext uses pdfcpu's context for more advanced extraction
func (e *pdfStringExtractor) ExtractStringsWithContext() ([]string, error) {
	// Read the PDF context
	ctx, err := api.ReadContextFile(e.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read PDF context: %w", err)
	}

	var allStrings []string

	// Extract strings from each page
	for pageNumber := 1; pageNumber <= ctx.PageCount; pageNumber++ {
		pageStrings, err := e.extractStringsFromPage(ctx, pageNumber)
		if err != nil {
			log.Printf("Warning: Failed to extract strings from page %d: %v", pageNumber, err)
			continue
		}
		allStrings = append(allStrings, pageStrings...)
	}

	return e.processStrings(allStrings), nil
}

// PrintStrings displays the extracted strings in a formatted way
func (e *pdfStringExtractor) PrintStrings(strings []string) {
	fmt.Printf("Found %d unique string literals in PDF: %s\n\n", len(strings), e.FilePath)

	for i, str := range strings {
		fmt.Printf("%3d. %s\n", i+1, str)

		// Print hex representation for debugging
		if len(str) < 50 { // Only for shorter strings
			fmt.Printf("     Hex: ")
			for _, r := range str {
				fmt.Printf("%02X ", r)
			}
			fmt.Println()
		}
		fmt.Println()
	}
}
