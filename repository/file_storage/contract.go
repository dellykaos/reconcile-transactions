package filestorage

type FileStorageRepository interface {
	Store(file *File) (string, error)
}
