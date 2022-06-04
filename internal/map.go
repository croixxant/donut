package internal

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

type Map struct {
	Src  File
	Dest File
}

type File struct {
	Path     string
	NotExist bool
	Lstat    fs.FileInfo
}

func NewFile(path string) (*File, error) {
	var notExist bool
	f, err := os.Lstat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		notExist = true
	}
	return &File{
		Path:     path,
		NotExist: notExist,
		Lstat:    f,
	}, nil
}

func (f *File) IsSymLink() bool {
	return f.Lstat.Mode()&os.ModeSymlink != 0
}

func (f *File) IsSameLink(path string) (bool, error) {
	if !f.IsSymLink() {
		return false, fmt.Errorf("%s: %w", f.Path, errors.New("not symlink"))
	}
	l, err := os.Readlink(f.Path)
	if err != nil {
		return false, fmt.Errorf("%s: %w", f.Path, err)
	}
	return l == path, nil
}
