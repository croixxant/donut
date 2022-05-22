package internal

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func Copy(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	return err
}

func Mkdir(dir string) error {
	if err := IsDir(dir); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

// IsDir checks SrcDir is exists and is directory
func IsDir(s string) error {
	if s == "" {
		return errors.New("not defined")
	} else if fi, err := os.Stat(s); err != nil {
		return err
	} else if !fi.IsDir() {
		return fmt.Errorf("%s is not directory", s)
	}
	return nil
}

// if dir is not directory, this returns home directory
func DirOrHome(dir string) (string, error) {
	if err := IsDir(dir); err != nil {
		return os.UserHomeDir()
	}
	return dir, nil
}

// Abs returns if path is relative, joins baseDir. if path is absolute, clean the path.
func Abs(path, baseDir string) string {
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}
	return filepath.Join(baseDir, path)
}
