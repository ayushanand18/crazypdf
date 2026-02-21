package crazypdf

import "errors"

var (
	// ErrInvalidPDF indicates the file is not a valid PDF or is corrupted.
	ErrInvalidPDF = errors.New("crazypdf: invalid or corrupted PDF")

	// ErrPasswordRequired indicates the PDF is encrypted and a password is needed.
	ErrPasswordRequired = errors.New("crazypdf: PDF is encrypted; password required")

	// ErrWrongPassword indicates the provided password is incorrect.
	ErrWrongPassword = errors.New("crazypdf: incorrect password")

	// ErrPageOutOfRange indicates the requested page index is out of bounds.
	ErrPageOutOfRange = errors.New("crazypdf: page index out of range")

	// ErrDocumentClosed indicates an operation was attempted on a closed document.
	ErrDocumentClosed = errors.New("crazypdf: document is closed")
)
