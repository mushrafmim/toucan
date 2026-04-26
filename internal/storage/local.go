package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type LocalStore struct {
	root    string
	baseURL string
}

func NewLocalStore(root, baseURL string) (*LocalStore, error) {
	if err := os.MkdirAll(root, 0755); err != nil {
		return nil, fmt.Errorf("create storage root: %w", err)
	}
	return &LocalStore{root: root, baseURL: baseURL}, nil
}

func (s *LocalStore) Upload(ctx context.Context, key string, data io.Reader) error {
	path := filepath.Join(s.root, key)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("create directories: %w", err)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, data); err != nil {
		return fmt.Errorf("copy data: %w", err)
	}

	return nil
}

func (s *LocalStore) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	path := filepath.Join(s.root, key)
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("open file: %w", err)
	}
	return f, nil
}

func (s *LocalStore) Delete(ctx context.Context, key string) error {
	path := filepath.Join(s.root, key)
	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return ErrNotFound
		}
		return fmt.Errorf("remove file: %w", err)
	}
	return nil
}

func (s *LocalStore) PresignUpload(ctx context.Context, key string, expires time.Duration) (string, error) {
	// For local storage, we just return the direct URL.
	// In a production app with local storage, you might generate a signed JWT.
	return fmt.Sprintf("%s/%s", s.baseURL, key), nil
}

func (s *LocalStore) PresignDownload(ctx context.Context, key string, expires time.Duration) (string, error) {
	return fmt.Sprintf("%s/%s", s.baseURL, key), nil
}
