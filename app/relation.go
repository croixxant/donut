package app

import (
	"io/fs"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slices"
)

type relation struct {
	Src  file
	Dest file
}

func newRelation(src, dest string) (*relation, error) {
	var sFile, dFile *file
	var err error
	if sFile, err = newFile(src); err != nil {
		return nil, err
	}
	if dFile, err = newFile(dest); err != nil {
		return nil, err
	}

	return &relation{
		Src:  *sFile,
		Dest: *dFile,
	}, nil
}

type relationsBuilder struct {
	srcDir   string
	destDir  string
	excludes []string
	mappings map[string]string
}

type relationsBuilderOption func(b *relationsBuilder)

func newRelationsBuilder(srcDir, destDir string, funcs ...relationsBuilderOption) *relationsBuilder {
	b := &relationsBuilder{
		srcDir:   srcDir,
		destDir:  destDir,
		excludes: []string{".git"},
		mappings: map[string]string{},
	}

	for _, fn := range funcs {
		fn(b)
	}

	return b
}

func withExcludes(s ...string) relationsBuilderOption {
	return func(b *relationsBuilder) {
		b.excludes = append(b.excludes, s...)
	}
}

func withMapping(m map[string]string) relationsBuilderOption {
	return func(b *relationsBuilder) {
		for k, v := range m {
			b.mappings[k] = v
		}
	}
}

func (b *relationsBuilder) build() ([]relation, error) {
	var rels []relation
	err := filepath.WalkDir(b.srcDir, func(path string, d fs.DirEntry, _ error) error {
		prefixTrimmed := strings.TrimPrefix(strings.TrimPrefix(path, b.srcDir), string(filepath.Separator))

		eq := func(s string) bool {
			ok, _ := filepath.Match(s, prefixTrimmed)
			return ok
		}
		if d.IsDir() { // skip directoires
			if slices.IndexFunc(b.excludes, eq) != -1 {
				return fs.SkipDir
			}
			return nil
		}

		// ignore b.excludes
		if slices.IndexFunc(b.excludes, eq) != -1 {
			return nil
		}

		// define destination path
		dPath := strings.Replace(path, b.srcDir, b.destDir, 1)
		// remap configured destination path
		if re := b.mappings[path]; re != "" {
			dPath = re
		}

		m, err := newRelation(path, dPath)
		if err != nil {
			return err
		}
		rels = append(rels, *m)
		return nil
	})

	return rels, err
}
