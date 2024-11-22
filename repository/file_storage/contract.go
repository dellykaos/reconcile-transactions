package filestorage

// FileStorageRepository is contract to store and get file
type FileStorageRepository interface {
	Store(file *File) (string, error)
	Get(filePath string) (*File, error)
}
