package storage

import (
	"context"
	"errors"
	"io"
	"time"
)

var (
	ErrNotFound = errors.New("file not found")
)

// Store defines the interface for blob storage operations.
type Store interface {
	// Upload stores the data with the given key.
	Upload(ctx context.Context, key string, data io.Reader) error

	// Download retrieves the data for the given key.
	Download(ctx context.Context, key string) (io.ReadCloser, error)

	// Delete removes the data for the given key.
	Delete(ctx context.Context, key string) error

	// PresignUpload generates a URL that can be used to upload data via PUT.
	PresignUpload(ctx context.Context, key string, expires time.Duration) (string, error)

	// PresignDownload generates a URL that can be used to download data via GET.
	PresignDownload(ctx context.Context, key string, expires time.Duration) (string, error)
}
