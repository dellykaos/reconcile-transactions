package localfilestorage

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"

	filestorage "github.com/delly/amartha/repository/file_storage"
)

// Storage is a struct to hold storage path
type Storage struct {
	storagePath string
}

var _ = filestorage.FileStorageRepository(&Storage{})

// NewStorage is a function to create new storage
func NewStorage(storagePath string) *Storage {
	return &Storage{
		storagePath: storagePath,
	}
}

// Store is a function to store file to local storage
func (lfs *Storage) Store(_ context.Context, file *filestorage.File) (string, error) {
	targetDir := filepath.Join(lfs.storagePath, file.Dir)
	if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create storage directory: %w", err)
	}

	targetPath := filepath.Join(targetDir, file.Name)
	outFile, err := os.Create(targetPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer outFile.Close()

	if _, err := outFile.Write(file.Buf.Bytes()); err != nil {
		return "", fmt.Errorf("failed to write file content: %w", err)
	}

	return targetPath, nil
}

// Get is a function to get file from local storage
func (lfs *Storage) Get(_ context.Context, filePath string) (*filestorage.File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file stat: %w", err)
	}

	buf := make([]byte, stat.Size())
	if _, err := file.Read(buf); err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	return &filestorage.File{
		Name: stat.Name(),
		Buf:  bytes.NewBuffer(buf),
	}, nil
}
