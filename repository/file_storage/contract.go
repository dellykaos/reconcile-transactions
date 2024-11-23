package filestorage

import "context"

// FileStorageRepository is contract to store and get file
type FileStorageRepository interface {
	Get(ctx context.Context, filePath string) (*File, error)
	Store(ctx context.Context, file *File) (string, error)
}
