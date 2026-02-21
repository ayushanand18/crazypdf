package crazypdf

// Config holds configuration for opening a PDF document.
type Config struct {
	// Password is the password for encrypted PDFs. Empty string for unencrypted.
	Password string
}

// Option is a functional option for configuring PDF document opening.
type Option func(*Config)

// WithPassword sets the password for opening encrypted PDFs.
func WithPassword(password string) Option {
	return func(c *Config) {
		c.Password = password
	}
}

// applyOptions creates a Config from the given options.
func applyOptions(opts []Option) *Config {
	cfg := &Config{}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}
