package internal

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

type Map struct {
	Src  string
	Dest string
}

var ErrAlreadyExists = errors.New("already exists")
var ErrAlreadyLinked = errors.New("already linked")

func (m *Map) CanLink() error {
	f, err := os.Lstat(m.Dest)
	if err != nil {
		if os.IsNotExist(err) { // if Lstat() returns not exists error
			return nil
		}
		return fmt.Errorf("%s: %w", m.Dest, err) // if Lstat() returns other error
	}

	if f.Mode()&os.ModeSymlink == 0 { // if not symlink
		return fmt.Errorf("%s: %w", m.Dest, fs.ErrExist)
	}

	link, err := os.Readlink(m.Dest)
	if err != nil { // if Readlink() returns error
		return fmt.Errorf("%s: %w", m.Dest, err)
	}
	if link != m.Src { // if link is not same as source path
		return fmt.Errorf("%s: %s", m.Dest, fs.ErrExist)
	}
	return ErrAlreadyLinked
}

func (m *Map) CanCopy() error {
	f, err := os.Lstat(m.Dest)
	if err != nil {
		if os.IsNotExist(err) { // if Lstat() returns not exists error
			return nil
		}
		return fmt.Errorf("%s: %w", m.Dest, err) // if Lstat() returns other error
	}

	if f.Mode()&os.ModeSymlink != 0 { // if symlink
		return fmt.Errorf("%s: %w", m.Dest, ErrAlreadyLinked)
	}

	return nil // if exists, returns noerror
}
