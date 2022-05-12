package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slices"

	"github.com/croixxant/donut/internal"
)

var (
	ignoreNames = []string{internal.FileMapConfigName, ".git"}
)

func findSourceDir(cfg *internal.Config) (string, error) {
	if cfg == nil {
		return "", errors.New("configuration file not found")
	}
	path := os.ExpandEnv(cfg.SrcDir)
	if path == "" {
		return "", errors.New("source_dir is not defined")
	} else if fi, err := os.Stat(path); err != nil {
		return "", err
	} else if !fi.IsDir() {
		return "", fmt.Errorf("%s is not directory", cfg.SrcDir)
	}
	return path, nil
}

func findDestinationDir(fileMapConfig *internal.FileMapConfig) (string, error) {
	if fileMapConfig == nil {
		return os.UserHomeDir()
	}
	path := os.ExpandEnv(fileMapConfig.DestDir)
	if path == "" {
		return os.UserHomeDir()
	} else if fi, err := os.Stat(path); err != nil {
		return os.UserHomeDir()
	} else if !fi.IsDir() {
		return os.UserHomeDir()
	}
	return path, nil
}

func newFileMaps(srcDir, destDir string, excludeNames []string) (list []internal.FileMap, err error) {
	err = filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		eq := func(s string) bool {
			return strings.Contains(strings.TrimPrefix(path, srcDir), s)
		}
		excludes := append(ignoreNames, excludeNames...)
		if slices.IndexFunc(excludes, eq) != -1 {
			return nil
		}
		list = append(list, internal.FileMap{
			Src:  path,
			Dest: strings.Replace(path, srcDir, destDir, 1),
		})
		return nil
	})
	return
}
