package strings

import (
	"context"
)

// PDFStringExtractor handles PDF string extraction operations
type pdfStringExtractor struct {
	FilePath string
	Context  context.Context
}

func NewPDFStringExtractor(filePath string) IPDFStringExtractor {
	return &pdfStringExtractor{
		FilePath: filePath,
		Context:  context.Background(),
	}
}

type IPDFStringExtractor interface {
	ExtractAllStrings() ([]string, error)
	ExtractStringsWithContext() ([]string, error)
	PrintStrings([]string)
}
