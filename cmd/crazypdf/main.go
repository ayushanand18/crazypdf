// Command crazypdf is a CLI tool for PDF processing.
//
// Usage:
//
//	crazypdf <command> [options] <input.pdf> [output]
//
// Commands:
//
//	text       Extract text from PDF
//
// Use "crazypdf <command> -h" for help on a specific command.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ayushanand18/crazypdf/pkg/crazypdf"
	"github.com/ayushanand18/crazypdf/pkg/extract"
)

const usage = `crazypdf - A PDF processing toolkit

Usage:
  crazypdf <command> [options] <input.pdf> [output.txt]

Commands:
  text       Extract text from a PDF file

Options vary by command. Use "crazypdf <command> -h" for help.

Examples:
  crazypdf text document.pdf
  crazypdf text -layout document.pdf output.txt
  crazypdf text -raw -pages 1-3 document.pdf
  crazypdf text -password secret encrypted.pdf
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "text":
		runTextCommand(os.Args[2:])
	case "-h", "--help", "help":
		fmt.Print(usage)
	case "-v", "--version", "version":
		fmt.Println("crazypdf v0.1.0")
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
		fmt.Fprint(os.Stderr, usage)
		os.Exit(1)
	}
}

func runTextCommand(args []string) {
	fs := flag.NewFlagSet("text", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, `Extract text from a PDF file.

Usage:
  crazypdf text [options] <input.pdf> [output.txt]

Options:
`)
		fs.PrintDefaults()
		fmt.Fprintf(os.Stderr, `
Examples:
  crazypdf text document.pdf
  crazypdf text -layout document.pdf output.txt
  crazypdf text -raw document.pdf
  crazypdf text -password secret encrypted.pdf
  crazypdf text -pages 1,3,5 document.pdf
`)
	}

	layout := fs.Bool("layout", false, "Preserve physical layout of text")
	raw := fs.Bool("raw", false, "Extract text in content stream order")
	password := fs.String("password", "", "Password for encrypted PDF")
	pagesFlag := fs.String("pages", "", "Page range (e.g., '1-5' or '1,3,5')")

	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}

	remaining := fs.Args()
	if len(remaining) < 1 {
		fmt.Fprintln(os.Stderr, "Error: input PDF file is required")
		fs.Usage()
		os.Exit(1)
	}

	inputFile := remaining[0]
	var outputFile string
	if len(remaining) > 1 {
		outputFile = remaining[1]
	}

	// determine layout mode
	var layoutMode extract.LayoutMode
	switch {
	case *layout:
		layoutMode = extract.LayoutPhysical
	case *raw:
		layoutMode = extract.LayoutRaw
	default:
		layoutMode = extract.LayoutSimple
	}

	// open document
	var docOpts []crazypdf.Option
	if *password != "" {
		docOpts = append(docOpts, crazypdf.WithPassword(*password))
	}

	doc, err := crazypdf.Open(inputFile, docOpts...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening PDF: %v\n", err)
		os.Exit(1)
	}
	defer doc.Close()

	// parse page range
	pageIndices, err := parsePageRange(*pagesFlag, doc.NumPages())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing page range: %v\n", err)
		os.Exit(1)
	}

	// extract text
	extractOpts := []extract.Option{
		extract.WithLayout(layoutMode),
	}

	var result strings.Builder
	for i, pageIdx := range pageIndices {
		page, err := doc.Page(pageIdx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accessing page %d: %v\n", pageIdx+1, err)
			os.Exit(1)
		}

		text, err := extract.PageText(page, extractOpts...)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error extracting text from page %d: %v\n", pageIdx+1, err)
			os.Exit(1)
		}

		result.WriteString(text)
		if i < len(pageIndices)-1 {
			result.WriteString("\n\n")
		}
	}

	output := result.String()
	if outputFile != "" {
		if err := os.WriteFile(outputFile, []byte(output), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Text extracted to %s (%d pages)\n", outputFile, len(pageIndices))
	} else {
		fmt.Print(output)
		if !strings.HasSuffix(output, "\n") {
			fmt.Println()
		}
	}
}

// parsePageRange parses a page range string like "1-5" or "1,3,5" into
// 0-based page indices.
func parsePageRange(pagesStr string, totalPages int) ([]int, error) {
	if pagesStr == "" {
		// All pages
		indices := make([]int, totalPages)
		for i := range indices {
			indices[i] = i
		}
		return indices, nil
	}

	var indices []int
	parts := strings.Split(pagesStr, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if strings.Contains(part, "-") {
			// Range: "1-5"
			var start, end int
			if _, err := fmt.Sscanf(part, "%d-%d", &start, &end); err != nil {
				return nil, fmt.Errorf("invalid page range: %q", part)
			}
			if start < 1 || end < 1 || start > totalPages || end > totalPages {
				return nil, fmt.Errorf("page range out of bounds: %q (document has %d pages)", part, totalPages)
			}
			if start > end {
				return nil, fmt.Errorf("invalid page range: start > end in %q", part)
			}
			for i := start; i <= end; i++ {
				indices = append(indices, i-1) // convert to 0-based
			}
		} else {
			// Single page: "3"
			var page int
			if _, err := fmt.Sscanf(part, "%d", &page); err != nil {
				return nil, fmt.Errorf("invalid page number: %q", part)
			}
			if page < 1 || page > totalPages {
				return nil, fmt.Errorf("page %d out of range (document has %d pages)", page, totalPages)
			}
			indices = append(indices, page-1) // convert to 0-based
		}
	}

	return indices, nil
}
