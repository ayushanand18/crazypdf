package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	crazystrings "github.com/ayushanand18/crazypdf/pkg/strings"
)

// extract strings from a PDF file path
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <pdf-file>")
		fmt.Println("Example: go run main.go document.pdf")
		os.Exit(1)
	}

	pdfFile := os.Args[1]

	// Check if file exists
	if _, err := os.Stat(pdfFile); os.IsNotExist(err) {
		log.Fatalf("PDF file not found: %s", pdfFile)
	}

	// Create extractor instance
	extractor := crazystrings.NewPDFStringExtractor(pdfFile)

	fmt.Printf("Analyzing PDF: %s\n", pdfFile)
	fmt.Println("=" + strings.Repeat("=", len(pdfFile)+12))

	// Method 1: Raw byte extraction (fast, but less accurate)
	fmt.Println("\n1. Raw byte extraction:")
	rawStrings, err := extractor.ExtractAllStrings()
	if err != nil {
		log.Fatalf("Error extracting strings: %v", err)
	}
	extractor.PrintStrings(rawStrings)

	// Method 2: Context-based extraction (slower, but more accurate)
	fmt.Println("\n2. Context-based extraction (using pdfcpu):")
	contextStrings, err := extractor.ExtractStringsWithContext()
	if err != nil {
		log.Printf("Warning: Context-based extraction failed: %v", err)
		log.Println("Falling back to raw extraction results")
	} else {
		extractor.PrintStrings(contextStrings)
	}

	// Summary
	fmt.Printf("Summary:\n")
	fmt.Printf("- Total unique strings found: %d\n", len(rawStrings))
	fmt.Printf("- PDF page count: (context-based extraction would show this)\n")
}
