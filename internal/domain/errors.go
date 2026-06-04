package domain

import "errors"

var (
	ErrNotPDF       = errors.New("file is not a PDF")
	ErrFileTooLarge = errors.New("file exceeds maximum size")
)
