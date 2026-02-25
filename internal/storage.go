package internal

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

const (
	PageSize = 8 * 1024
	DataDir  = "/var/lib/evergreen"
)

type Page struct {
	buf [PageSize]byte
}

func NewPage() *Page {
	return &Page{}
}

func (p *Page) Bytes() []byte {
	return p.buf[:]
}

type Storage interface {
	ReadPage(p *Page, pageNum int64) error
	WritePage(p *Page, pageNum int64) error
}

type FileStore struct {
	f *os.File
}

var ErrIndexAlreadyExists = errors.New("index already exists")

func NewFileStore(name string) error {
	path := filepath.Join(DataDir, name)

	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("%w: %s", ErrIndexAlreadyExists, name)
	} else if !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("failed to check existence of index: %w", err)
	}

	_, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create index file: %w", err)
	}

	return nil
}

func (fs *FileStore) ReadPage(p *Page, pageNum int64) error {
	off := PageSize * pageNum

	n, err := fs.f.ReadAt(p.Bytes(), off)
	if err != nil {
		return err
	}
	if n != PageSize {
		return fmt.Errorf("short read: %d", n)
	}

	return nil
}

func (fs *FileStore) WritePage(p *Page, pageNum int64) error {
	if len(p.Bytes()) != PageSize {
		return fmt.Errorf("invalid page size: got %d bytes", len(p.Bytes()))
	}

	off := PageSize * pageNum
	n, err := fs.f.WriteAt(p.Bytes(), off)
	if err != nil {
		return err
	}

	if n != PageSize {
		return fmt.Errorf("short write: expected %d bytes, got %d", PageSize, n)
	}

	return nil
}
