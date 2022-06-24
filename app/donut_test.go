package app

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/croixxant/donut/testutil"
)

func TestNew(t *testing.T) {
	dir := t.TempDir()

	type args struct {
		testdata map[string]interface{}
	}
	tests := []struct {
		name      string
		args      args
		want      *Donut
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "OK",
			args: args{
				testdata: map[string]interface{}{
					"src_dir":  dir,
					"dest_dir": dir,
					"excludes": []string{"example.toml"},
					"method":   "link",
					"mappings": map[string]string{"foo": "bar"},
				},
			},
			want: &Donut{
				config: &Config{
					SrcDir:   dir,
					DestDir:  dir,
					Excludes: []string{"example.toml"},
					Method:   "link",
					Mappings: map[string]string{"foo": "bar"},
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "Error",
			args: args{
				testdata: map[string]interface{}{
					"src_dir":  "/home/gopher",
					"dest_dir": "/home/gopher",
				},
			},
			want: &Donut{
				config: &Config{},
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := viper.New()
			if err := v.MergeConfigMap(tt.args.testdata); err != nil {
				t.Fatal(err)
			}
			o, e := new(bytes.Buffer), new(bytes.Buffer)
			tt.want.out, tt.want.err = o, e
			got, err := New(WithConfig(v), WithOut(o), WithErr(e))
			tt.assertion(t, err)
			if err == nil {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestDonut_Init(t *testing.T) {
	testutil.SetDirEnv(t, UserHomeDir)

	type fields struct {
		config *Config
	}
	type args struct {
		srcDir  string
		cfgPath string
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		defaultConfig string
		want          map[string]interface{}
		assertion     assert.ErrorAssertionFunc
	}{
		{
			name: "OK",
			fields: fields{
				config: nil,
			},
			args: args{
				srcDir:  "",
				cfgPath: "",
			},
			want: map[string]interface{}{
				"src_dir":  filepath.Join(UserHomeDir, ".local", "share", Name),
				"dest_dir": UserHomeDir,
				"method":   MethodCopy,
			},
			assertion: assert.NoError,
		},
		{
			name: "OK/SrcDirSpecified",
			fields: fields{
				config: nil,
			},
			args: args{
				srcDir:  filepath.Join(UserHomeDir, Name),
				cfgPath: "",
			},
			want: map[string]interface{}{
				"src_dir":  filepath.Join(UserHomeDir, Name),
				"dest_dir": UserHomeDir,
				"method":   MethodCopy,
			},
			assertion: assert.NoError,
		},
		{
			name: "OK/CfgPathSpecified",
			fields: fields{
				config: nil,
			},
			args: args{
				srcDir:  "",
				cfgPath: filepath.Join(UserHomeDir, ".config", Name+".toml"),
			},
			want: map[string]interface{}{
				"src_dir":  filepath.Join(UserHomeDir, ".local", "share", Name),
				"dest_dir": UserHomeDir,
				"method":   MethodCopy,
			},
			assertion: assert.NoError,
		},
		{
			name: "Error/NoNeed",
			fields: fields{
				config: nil,
			},
			args: args{
				srcDir:  "",
				cfgPath: "",
			},
			defaultConfig: "../testdata/config/basic.toml",
			want:          map[string]interface{}{},
			assertion:     assert.Error,
		},
		{
			name: "Error/FileBroken",
			fields: fields{
				config: nil,
			},
			args: args{
				srcDir:  "",
				cfgPath: "",
			},
			defaultConfig: "../testdata/config/broken.toml",
			want:          map[string]interface{}{},
			assertion:     assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.defaultConfig != "" {
				defer testutil.CopyFile(t, tt.defaultConfig, filepath.Join(UserHomeDir, ".config", Name+".toml"))()
			}
			o, e := new(bytes.Buffer), new(bytes.Buffer)
			d, _ := New(WithOut(o), WithErr(e))
			cfgPath := filepath.Join(UserHomeDir, ".config", Name, Name+".toml")
			if tt.args.cfgPath != "" {
				cfgPath = tt.args.cfgPath
			}
			err := d.Init(tt.args.srcDir, tt.args.cfgPath)
			defer func() { os.Remove(cfgPath) }()
			tt.assertion(t, err)
			if err != nil {
				return
			}
			if diff := cmp.Diff(tt.want, testutil.ReadToml(t, cfgPath)); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDonut_Where(t *testing.T) {
	dir := t.TempDir()

	type fields struct {
		config *Config
	}
	tests := []struct {
		name      string
		fields    fields
		want      string
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "OK",
			fields: fields{
				config: &Config{
					SrcDir: dir,
				},
			},
			want:      dir + "\n",
			assertion: assert.NoError,
		},
		{
			name: "Error/NoConfig",
			fields: fields{
				config: nil,
			},
			want:      "",
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o, e := new(bytes.Buffer), new(bytes.Buffer)
			d, _ := New(WithOut(o), WithErr(e))
			d.config = tt.fields.config
			tt.assertion(t, d.Where())
			assert.Equal(t, tt.want, o.String())
		})
	}
}

func TestDonut_List(t *testing.T) {
	srcFiles := map[string]string{
		"../testdata/data/.example":     ".example",
		"../testdata/data/example.toml": "example.toml",
		"../testdata/data/example.yaml": ".config/example.yaml",
		"../testdata/data/example.json": "example.json",
		"../testdata/data/example.ini":  "example.ini",
		"../testdata/data/.zshrc":       ".zshrc",
		"../testdata/data/.zprofile":    ".zprofile",
		"../testdata/data/.zshenv":      ".zshenv",
		"../testdata/data/.zlogin":      ".zlogin",
		"../testdata/data/.empty":       ".git/.empty",
		"../testdata/data/.gitconfig":   ".gitconfig",
	}
	destFiles := map[string]string{"../testdata/data/.example": ".example"}
	links := map[string]string{
		".zprofile": ".zprofile",
		".zshrc":    ".zlogin",
	}
	mappings := map[string]string{
		"example.toml": ".config/example.toml",
	}
	excludes := []string{"example.json"}

	type fields struct {
		config *Config
	}
	tests := []struct {
		name      string
		srcFiles  map[string]string // map[testdata]srcdata
		destFiles map[string]string // map[testdata]destdata
		links     map[string]string // map[srcdata]destdata
		fields    fields
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "OK/Copy",
			srcFiles:  srcFiles,
			destFiles: destFiles,
			links:     links,
			fields: fields{
				config: &Config{
					Excludes: excludes,
					Method:   MethodCopy,
					Mappings: mappings,
				},
			},
			assertion: assert.NoError,
		},
		{
			name:      "OK/Link",
			srcFiles:  srcFiles,
			destFiles: destFiles,
			links:     links,
			fields: fields{
				config: &Config{
					Excludes: excludes,
					Method:   MethodLink,
					Mappings: mappings,
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "Error/NoConfig",
			fields: fields{
				config: nil,
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			home, _, data := testutil.CreateBaseDir(t)
			testutil.CreateDirs(t, filepath.Join(data, ".config"), filepath.Join(data, ".git"))
			for s, d := range tt.srcFiles {
				defer testutil.CopyFile(t, s, filepath.Join(data, d))()
			}
			for s, d := range tt.destFiles {
				defer testutil.CopyFile(t, s, filepath.Join(home, d))()
			}
			for s, d := range tt.links {
				if err := os.Symlink(filepath.Join(data, s), filepath.Join(home, d)); err != nil {
					t.Fatal(err)
				}
				defer os.Remove(d)
			}
			o, e := new(bytes.Buffer), new(bytes.Buffer)
			d, _ := New(WithOut(o), WithErr(e))
			if tt.fields.config != nil {
				d.config = tt.fields.config
				d.config.SrcDir = data
				d.config.DestDir = home
			}
			tt.assertion(t, d.List())
			// TODO: assert output
		})
	}
}

func TestDonut_Apply(t *testing.T) {
	srcFiles := map[string]string{
		"../testdata/data/.example":     ".example",
		"../testdata/data/example.toml": "example.toml",
		"../testdata/data/example.yaml": ".config/example.yaml",
		"../testdata/data/example.json": "example.json",
		"../testdata/data/example.ini":  "example.ini",
		"../testdata/data/.zshrc":       ".zshrc",
		"../testdata/data/.zprofile":    ".zprofile",
		"../testdata/data/.zshenv":      ".zshenv",
		"../testdata/data/.zlogin":      ".zlogin",
		"../testdata/data/.empty":       ".git/.empty",
		"../testdata/data/.gitconfig":   ".gitconfig",
	}
	destFiles := map[string]string{"../testdata/data/.example": ".example"}
	links := map[string]string{
		".zprofile": ".zprofile",
		".zshrc":    ".zlogin",
	}
	mappings := map[string]string{
		"example.toml": ".config/example.toml",
	}
	excludes := []string{"example.json"}

	type fields struct {
		config *Config
	}
	type args struct {
		force bool
	}
	tests := []struct {
		name      string
		srcFiles  map[string]string // map[testdata]srcdata
		destFiles map[string]string // map[testdata]destdata
		links     map[string]string // map[srcdata]destdata
		fields    fields
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "OK/Copy",
			srcFiles:  srcFiles,
			destFiles: destFiles,
			links:     links,
			fields: fields{
				config: &Config{
					Excludes: excludes,
					Method:   MethodCopy,
					Mappings: mappings,
				},
			},
			args: args{
				force: false,
			},
			assertion: assert.NoError,
		},
		{
			name:      "OK/Link",
			srcFiles:  srcFiles,
			destFiles: destFiles,
			links:     links,
			fields: fields{
				config: &Config{
					Excludes: excludes,
					Method:   MethodLink,
					Mappings: mappings,
				},
			},
			args: args{
				force: false,
			},
			assertion: assert.NoError,
		},
		{
			name:      "OK/Copy/Force",
			srcFiles:  srcFiles,
			destFiles: destFiles,
			links:     links,
			fields: fields{
				config: &Config{
					Excludes: excludes,
					Method:   MethodCopy,
					Mappings: mappings,
				},
			},
			args: args{
				force: true,
			},
			assertion: assert.NoError,
		},
		{
			name:      "OK/Link/Force",
			srcFiles:  srcFiles,
			destFiles: destFiles,
			links:     links,
			fields: fields{
				config: &Config{
					Excludes: excludes,
					Method:   MethodLink,
					Mappings: mappings,
				},
			},
			args: args{
				force: true,
			},
			assertion: assert.NoError,
		},
		{
			name: "Error/NoConfig",
			fields: fields{
				config: nil,
			},
			args: args{
				force: false,
			},
			assertion: assert.Error,
		},
		{
			name:      "Error/Link/OtherLinkExists",
			srcFiles:  srcFiles,
			destFiles: destFiles,
			links:     links,
			fields: fields{
				config: &Config{
					Excludes: excludes,
					Method:   MethodLink,
					Mappings: map[string]string{
						"example.toml": ".config/example.toml",
						".zshenv":      ".zshrc",
					},
				},
			},
			args: args{
				force: true,
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			home, _, data := testutil.CreateBaseDir(t)
			testutil.CreateDirs(t, filepath.Join(data, ".config"), filepath.Join(data, ".git"))
			for s, d := range tt.srcFiles {
				defer testutil.CopyFile(t, s, filepath.Join(data, d))()
			}
			for s, d := range tt.destFiles {
				defer testutil.CopyFile(t, s, filepath.Join(home, d))()
			}
			for s, d := range tt.links {
				if err := os.Symlink(filepath.Join(data, s), filepath.Join(home, d)); err != nil {
					t.Fatal(err)
				}
				defer os.Remove(d)
			}
			o, e := new(bytes.Buffer), new(bytes.Buffer)
			d, _ := New(WithOut(o), WithErr(e))
			if tt.fields.config != nil {
				d.config = tt.fields.config
				d.config.SrcDir = data
				d.config.DestDir = home
			}
			tt.assertion(t, d.Apply(tt.args.force))
			// TODO: assert output
		})
	}
}
