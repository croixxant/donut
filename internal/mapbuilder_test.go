package internal

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMapBuilder(t *testing.T) {
	type args struct {
		srcDir  string
		destDir string
		funcs   []MapBuilderOption
	}
	tests := []struct {
		name string
		args args
		want *MapBuilder
	}{
		{
			name: "OK",
			args: args{
				srcDir:  "/src",
				destDir: "/dest",
				funcs: []MapBuilderOption{
					WithExcludes("README.md"),
					WithRemaps(map[string]string{".zshenv": "/home/gopher/.zshenv"}),
				},
			},
			want: &MapBuilder{
				srcDir:   "/src",
				destDir:  "/dest",
				excludes: []string{"README.md"},
				remaps:   map[string]string{".zshenv": "/home/gopher/.zshenv"},
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewMapBuilder(tt.args.srcDir, tt.args.destDir, tt.args.funcs...))
		})
	}
}

func TestMapBuilder_Build(t *testing.T) {
	srcDir, destDir := t.TempDir(), t.TempDir()
	files := []string{"README.md", ".starship.toml", ".zshenv"}
	for _, name := range files {
		if err := os.WriteFile(filepath.Join(srcDir, name), []byte(""), os.ModePerm); err != nil {
			t.Fatalf("tmpFile not created: %s", err)
		}
	}
	want := []Map{}
	wantFiles := []struct {
		Src  string
		Dest string
	}{
		{".starship.toml", "/.config/starship.toml"},
		{".zshenv", ".zshenv"},
	}
	for _, v := range wantFiles {
		sFile, err := newFile(filepath.Join(srcDir, v.Src))
		if err != nil {
			t.Fatal(err)
		}
		dFile, err := newFile(filepath.Join(destDir, v.Dest))
		if err != nil {
			t.Fatal(err)
		}
		want = append(want, newMap(sFile, dFile))
	}

	type fields struct {
		srcDir   string
		destDir  string
		excludes []string
		remaps   map[string]string
	}
	tests := []struct {
		name      string
		fields    fields
		want      []Map
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "OK",
			fields: fields{
				srcDir:   srcDir,
				destDir:  destDir,
				excludes: []string{"README.md"},
				remaps: map[string]string{
					filepath.Join(srcDir, ".starship.toml"): filepath.Join(destDir, "/.config/starship.toml"),
				},
			},
			want:      want,
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &MapBuilder{
				srcDir:   tt.fields.srcDir,
				destDir:  tt.fields.destDir,
				excludes: tt.fields.excludes,
				remaps:   tt.fields.remaps,
			}
			got, err := b.Build()
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
