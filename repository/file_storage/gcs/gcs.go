package gcs

import (
	"bytes"
	"context"
	"io"
	"path/filepath"

	"cloud.google.com/go/storage"
	"github.com/delly/amartha/common/logger"
	filestorage "github.com/delly/amartha/repository/file_storage"
	"go.uber.org/zap"
)

// Bucket is a struct to store file in GCS
type Bucket struct {
	bucket *storage.BucketHandle
	log    *zap.Logger
}

// NewBucket create new GCS Bucket
func NewBucket(bucket *storage.BucketHandle) *Bucket {
	return &Bucket{
		bucket: bucket,
		log:    zap.L().With(zap.String("repository", "gcs")),
	}
}

// Get get file from GCS Bucket
func (b *Bucket) Get(ctx context.Context, filePath string) (*filestorage.File, error) {
	log := logger.WithMethod(b.log, "Get")
	obj := b.bucket.Object(filePath)
	r, err := obj.NewReader(ctx)
	if err != nil {
		log.Error("failed to read file", zap.Error(err), zap.String("file", filePath))
		return nil, err
	}

	buf, err := io.ReadAll(r)
	if err != nil {
		log.Error("failed to read file", zap.Error(err), zap.String("file", filePath))
		return nil, err
	}

	return &filestorage.File{
		Name: filepath.Base(filePath),
		Dir:  filepath.Dir(filePath),
		Buf:  bytes.NewBuffer(buf),
	}, nil
}

// Store store file to GCS Bucket
func (b *Bucket) Store(ctx context.Context, file *filestorage.File) (string, error) {
	log := logger.WithMethod(b.log, "Store")
	target := filepath.Join(file.Dir, file.Name)
	obj := b.bucket.Object(target)
	w := obj.NewWriter(ctx)

	if _, err := io.Copy(w, file.Buf); err != nil {
		log.Error("failed to write file", zap.Error(err), zap.String("file", target))
		return "", err
	}

	if err := w.Close(); err != nil {
		log.Error("failed to close writer", zap.Error(err), zap.String("file", target))
		return "", err
	}

	return target, nil
}
