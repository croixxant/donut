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
	remaps   []Map
}

type MapBuilderOption func(b *MapBuilder)

var ignores = []string{MapConfigName, ".git"}

func NewMapBuilder(srcDir, destDir string, funcs ...MapBuilderOption) *MapBuilder {
	b := &MapBuilder{
		srcDir:   srcDir,
		destDir:  destDir,
		excludes: ignores,
	}

	for _, fn := range funcs {
		fn(b)
	}
	return b
}

func WithExcludes(s []string) MapBuilderOption {
	return func(b *MapBuilder) {
		b.excludes = append(b.excludes, s...)
	}
}

func WithRemaps(s []Map) MapBuilderOption {
	return func(b *MapBuilder) {
		b.remaps = append(b.remaps, s...)
	}
}

func (b *MapBuilder) Build() []Map {
	remaps := map[string]string{}
	for _, m := range b.remaps {
		remaps[m.Src] = m.Dest
	}
	var maps []Map

	_ = filepath.WalkDir(b.srcDir, func(path string, d fs.DirEntry, err error) error {
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
		if re := remaps[path]; re != "" {
			dPath = re
		}

		maps = append(maps, Map{
			Src:  path,
			Dest: dPath,
		})
		return nil
	})

	return maps
}
