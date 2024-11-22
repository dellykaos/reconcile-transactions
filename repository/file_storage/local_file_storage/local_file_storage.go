package localfilestorage

import (
	"bytes"
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
	uniqueFileDirectory := lfs.generateUniqueFileDirectory()
	if err := os.MkdirAll(uniqueFileDirectory, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create storage directory: %w", err)
	}

	uniqueFilePath := filepath.Join(uniqueFileDirectory, file.Name)
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

// Get is a function to get file from local storage
func (lfs *Storage) Get(filePath string) (*filestorage.File, error) {
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

// TODO: Move generation unique path to service layer.
// Repository should not store any business logic.
func (lfs *Storage) generateUniqueFileDirectory() string {
	timestamp := time.Now().UnixNano()
	randomString := lfs.generateRandomString(8)
	uniqueFileDirectory := fmt.Sprintf("%d_%s", timestamp, randomString)
	return filepath.Join(lfs.storagePath, uniqueFileDirectory)
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
