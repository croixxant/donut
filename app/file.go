package app

import (
	"fmt"
	"io/fs"
	"os"
)

type file struct {
	Path     string
	NotExist bool
	Lstat    fs.FileInfo
}

func newFile(path string) (*file, error) {
	var notExist bool
	f, err := os.Lstat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		notExist = true
	}
	return &file{
		Path:     path,
		NotExist: notExist,
		Lstat:    f,
	}, nil
}

func (f *file) isSymLink() bool {
	return f.Lstat.Mode()&os.ModeSymlink != 0
}

func (f *file) isSame(path string) (bool, error) {
	if !f.isSymLink() {
		return f.Path == path, nil
	}
	l, err := os.Readlink(f.Path)
	if err != nil {
		return false, fmt.Errorf("%s: %w", f.Path, err)
	}
	return l == path, nil
}
