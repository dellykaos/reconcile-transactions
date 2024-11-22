package http

import (
	"fmt"
)

var (
	// ErrFileCannotBeAccessed is an error when file cannot be accessed
	ErrFileCannotBeAccessed = func(fname string) error {
		return fmt.Errorf("file %s cannot be accessed", fname)
	}
	// ErrExtensionFileUnknown is an error when file extension is unknown
	ErrExtensionFileUnknown = func(fname string) error {
		return fmt.Errorf("file %s must have a .csv extension", fname)
	}
	// ErrExtensionFileInvalid is an error when file extension is invalid
	ErrExtensionFileInvalid = func(fname string) error {
		return fmt.Errorf("file %s must have a .csv extension", fname)
	}
	// ErrFileSizeExceedLimit is an error when file size exceed limit
	ErrFileSizeExceedLimit = func(fname, limit string) error {
		return fmt.Errorf("file size %s more than %s", fname, limit)
	}
)
