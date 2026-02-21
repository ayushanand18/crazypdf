package extract

// LayoutMode controls how text is extracted from PDF pages.
type LayoutMode int

const (
	// LayoutSimple extracts plain text with minimal formatting.
	// Words are joined by spaces, rows by newlines.
	LayoutSimple LayoutMode = iota

	// LayoutRaw extracts text in content stream order, preserving
	// the order in which text appears in the PDF's internal structure.
	LayoutRaw

	// LayoutPhysical attempts to preserve the physical/spatial layout
	// of text on the page, using x,y coordinates to position text.
	LayoutPhysical
)

// textConfig holds configuration for text extraction operations.
type textConfig struct {
	Layout        LayoutMode
	PageSeparator string
	PageWidth     float64 // page width in points for physical layout
}

// Option is a functional option for configuring text extraction.
type Option func(*textConfig)

// WithLayout sets the text extraction layout mode.
func WithLayout(mode LayoutMode) Option {
	return func(c *textConfig) {
		c.Layout = mode
	}
}

// WithPageSeparator sets the separator string between pages when
// extracting text from the entire document.
func WithPageSeparator(sep string) Option {
	return func(c *textConfig) {
		c.PageSeparator = sep
	}
}

// WithPageWidth sets the page width in PDF points for physical layout mode.
// Default is 612 (US Letter width). This affects column spacing in LayoutPhysical.
func WithPageWidth(width float64) Option {
	return func(c *textConfig) {
		c.PageWidth = width
	}
}

// defaultConfig returns the default text extraction configuration.
func defaultConfig() *textConfig {
	return &textConfig{
		Layout:        LayoutSimple,
		PageSeparator: "\n\n",
		PageWidth:     612,
	}
}

// applyOptions creates a textConfig from the given options.
func applyOptions(opts []Option) *textConfig {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}
