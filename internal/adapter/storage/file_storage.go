package storage

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

func (s *Storage) UploadFile(src string) (_ string, err error) {
	f, err := os.Open(src)
	if err != nil {
		return "", fmt.Errorf("os open: %w", err)
	}
	defer func() {
		err = errors.Join(err, f.Close())
	}()

	filename := filepath.Join(s.filesDir(), uuid.NewString())

	d, err := os.Create(filename)
	if err != nil {
		return "", fmt.Errorf("os create: %w", err)
	}
	defer func() {
		err = errors.Join(err, d.Close())
	}()

	if _, err := io.Copy(d, f); err != nil {
		return "", fmt.Errorf("io copy: %w", err)
	}

	if err := d.Sync(); err != nil {
		return "", fmt.Errorf("sync: %w", err)
	}

	return filename, nil
}

func (s *Storage) DownloadFile(src, dst string) error {
	return nil
}
