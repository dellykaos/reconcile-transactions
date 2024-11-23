package http

import (
	"errors"
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
	// ErrBankTrxFileEmpty is an error when bank transaction files is empty
	ErrBankTrxFileEmpty = errors.New("bank transaction files is required, at least provide one")
	// ErrBankFileAndNameLengthNotMatch is an error when bank names and bank transaction files length not match
	ErrBankFileAndNameLengthNotMatch = errors.New("bank names and bank transaction files length must be same")
)
