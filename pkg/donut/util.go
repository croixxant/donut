package donut

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

var UserHomeDir string

func init() {
	var err error
	if UserHomeDir, err = os.UserHomeDir(); err != nil {
		panic(err)
	}
}

func SetUserHomeDir(dir string) func() {
	tmp := UserHomeDir
	UserHomeDir = dir
	return func() { UserHomeDir = tmp }
}

func DefaultConfigFile() string {
	return filepath.Join(UserHomeDir, ".config", Name, Name+".toml")
}

func DefaultConfigDirs() []string {
	return []string{
		"$XDG_CONFIG_HOME",
		filepath.Join("$XDG_CONFIG_HOME", Name),
		filepath.Join(UserHomeDir, ".config"),
		filepath.Join(UserHomeDir, ".config", Name),
	}
}

func DefaultSourceDir() string {
	return filepath.Join(UserHomeDir, ".local", "share", Name)
}

// absPath returns if path is relative, joins baseDir. if path is absolute, clean the path.
func absPath(path, baseDir string) string {
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}
	return filepath.Join(baseDir, path)
}

// isDir checks SrcDir is exists and is directory
func isDir(s string) error {
	if s == "" {
		return errors.New("not defined")
	} else if f, err := os.Stat(s); err != nil {
		return err
	} else if !f.IsDir() {
		return errors.New("not directory")
	}
	return nil
}

func copyFile(src, dest string) error {
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
