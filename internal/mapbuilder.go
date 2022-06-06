package internal

import (
	"io/fs"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slices"
)

type MapBuilder struct {
	srcDir   string
	destDir  string
	excludes []string
	remaps   map[string]string
}

type MapBuilderOption func(b *MapBuilder)

func NewMapBuilder(srcDir, destDir string, funcs ...MapBuilderOption) *MapBuilder {
	b := &MapBuilder{
		srcDir:   srcDir,
		destDir:  destDir,
		excludes: []string{},
		remaps:   map[string]string{},
	}

	for _, fn := range funcs {
		fn(b)
	}
	return b
}

func WithExcludes(s ...string) MapBuilderOption {
	return func(b *MapBuilder) {
		b.excludes = append(b.excludes, s...)
	}
}

func WithRemaps(m map[string]string) MapBuilderOption {
	return func(b *MapBuilder) {
		for k, v := range m {
			b.remaps[k] = v
		}
	}
}

func (b *MapBuilder) Build() ([]Map, error) {
	var maps []Map
	err := filepath.WalkDir(b.srcDir, func(path string, d fs.DirEntry, _ error) error {
		if d.IsDir() { // skip directoires
			return nil
		}

		// ignore b.excludes
		eq := func(s string) bool {
			return strings.Contains(strings.TrimPrefix(path, b.srcDir), s)
		}
		if slices.IndexFunc(b.excludes, eq) != -1 {
			return nil
		}

		// define destination path
		dPath := strings.Replace(path, b.srcDir, b.destDir, 1)
		// remap configured destination path
		if re := b.remaps[path]; re != "" {
			dPath = re
		}

		sFile, err := newFile(path)
		if err != nil {
			return err
		}
		dFile, err := newFile(dPath)
		if err != nil {
			return err
		}
		maps = append(maps, newMap(sFile, dFile))
		return nil
	})

	return maps, err
}
