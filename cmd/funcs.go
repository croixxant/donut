package cmd

import (
	"io/fs"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slices"

	"github.com/croixxant/donut/internal"
)

var (
	ignoreNames = []string{internal.MapConfigName, ".git"}
)

func newMap(srcDir, destDir string, mapConfig *internal.MapConfigData) (list []internal.Map, err error) {
	err = filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		// ignore default and configured source files
		eq := func(s string) bool {
			return strings.Contains(strings.TrimPrefix(path, srcDir), s)
		}
		excludes := append(ignoreNames, mapConfig.Excludes...)
		if slices.IndexFunc(excludes, eq) != -1 {
			return nil
		}

		// define destination path
		dPath := strings.Replace(path, srcDir, destDir, 1)

		// remap configured destination path
		m := mapConfig.Map(srcDir)
		if re := m[path]; re != "" {
			dPath = re
		}
		list = append(list, internal.Map{
			Src:  path,
			Dest: dPath,
		})
		return nil
	})
	return
}
