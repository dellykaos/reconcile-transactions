package filestorage

import "bytes"

// File is a struct to hold metadata of csv file
type File struct {
	Name string
	Dir  string
	Buf  *bytes.Buffer
}
