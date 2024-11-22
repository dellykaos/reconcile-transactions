package localfilestorage

import (
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"time"

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
func (lfs *Storage) Store(file *filestorage.File) (string, error) {
	if err := os.MkdirAll(lfs.storagePath, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create storage directory: %w", err)
	}

	uniqueFilePath := lfs.generateUniqueFilePath(file.Name)

	outFile, err := os.Create(uniqueFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer outFile.Close()

	if _, err := outFile.Write(file.Buf.Bytes()); err != nil {
		return "", fmt.Errorf("failed to write file content: %w", err)
	}

	return uniqueFilePath, nil
}

func (lfs *Storage) generateUniqueFilePath(fileName string) string {
	timestamp := time.Now().UnixNano()
	randomString := lfs.generateRandomString(8)
	uniqueFileName := fmt.Sprintf("%d_%s_%s", timestamp, randomString, fileName)
	return filepath.Join(lfs.storagePath, uniqueFileName)
}

func (lfs *Storage) generateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	for i := range bytes {
		bytes[i] = letters[bytes[i]%byte(len(letters))]
	}
	return string(bytes)
}
