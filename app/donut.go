package app

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/pterm/pterm"
	"github.com/spf13/viper"
)

const Name string = "donut"

type Donut struct {
	config *Config
	out    io.Writer
	err    io.Writer
}

func New(opts ...Option) (*Donut, error) {
	app := &Donut{
		out: os.Stdout,
		err: os.Stderr,
	}
	for _, opt := range opts {
		if err := opt(app); err != nil {
			return nil, err
		}
	}

	pterm.SetDefaultOutput(app.out)

	return app, nil
}

type Option func(*Donut) error

func WithConfig(v *viper.Viper) Option {
	return func(d *Donut) error {
		var cfg Config
		if err := v.Unmarshal(&cfg, viper.DecodeHook(defaultDecodeHookFunc)); err != nil {
			return err
		}
		if err := validateConfig(&cfg); err != nil {
			return err
		}
		d.config = &cfg
		return nil
	}
}

func WithOut(r io.Writer) Option {
	return func(d *Donut) error {
		d.out = r
		return nil
	}
}

func WithErr(r io.Writer) Option {
	return func(d *Donut) error {
		d.err = r
		return nil
	}
}

func validateConfig(d *Config) error {
	validateIsDir := func(s string) error {
		if err := isDir(s); err != nil {
			return err
		}
		return nil
	}

	for _, s := range []string{d.SrcDir, d.DestDir} {
		if err := validateIsDir(s); err != nil {
			return err
		}
	}

	return nil
}

func (d *Donut) Init(srcDir, cfgPath string) (err error) {
	var v *viper.Viper
	if cfgPath == "" {
		_, err = NewConfig(WithDefault(), WithNameAndPath(Name, DefaultConfigDirs()...))
		cfgPath = DefaultConfigFile()
	} else {
		_, err = NewConfig(WithFile(cfgPath))
	}
	if err == nil {
		return errors.New("no need to init")
	}
	if !errors.As(err, &viper.ConfigFileNotFoundError{}) && !os.IsNotExist(err) {
		return fmt.Errorf("config file already exists, but error: %w", err)
	}

	v, _ = NewConfig(WithDefault())

	if srcDir != "" {
		v.Set("src_dir", srcDir)
	}
	if err := os.MkdirAll(filepath.Dir(cfgPath), os.ModePerm); err != nil {
		return err
	}
	if err := WriteConfig(cfgPath); err != nil {
		return err
	}

	pterm.Success.Printfln("Configuration file created in %s", cfgPath)
	return nil
}

func (d *Donut) Where() error {
	if d.config == nil {
		return errors.New("no config")
	}
	fmt.Fprintln(d.out, d.config.SrcDir)
	return nil
}

func (d *Donut) List() error {
	if d.config == nil {
		return errors.New("no config")
	}

	list, err := newRelationsBuilder(
		d.config.SrcDir,
		d.config.DestDir,
		withExcludes(d.config.Excludes...),
		withMapping(d.config.absMappings()),
	).build()
	if err != nil {
		return err
	}

	tableData := make([][]string, 0, len(list)+1) // add header capacity
	header := []string{"SOURCE", "DESTINATION"}
	if d.config.Method == MethodLink {
		header = append([]string{"✔ "}, header...)
	}
	tableData = append(tableData, header)
	for _, v := range list {
		row := []string{v.Src.Path, v.Dest.Path}
		if d.config.Method == MethodLink {
			var l string
			if !v.Dest.NotExist {
				if linked, err := v.Dest.isSame(v.Src.Path); err != nil {
					return err
				} else if linked {
					l = "✔ "
				}
			}
			row = append([]string{l}, row...)
		}
		tableData = append(tableData, row)
	}

	if err := pterm.DefaultTable.
		WithHasHeader().
		WithData(tableData).
		WithBoxed().
		Render(); err != nil {
		return err
	}
	return nil
}

func (d *Donut) Apply(force bool) error {
	if d.config == nil {
		return errors.New("no config")
	}

	list, err := newRelationsBuilder(
		d.config.SrcDir,
		d.config.DestDir,
		withExcludes(d.config.Excludes...),
		withMapping(d.config.absMappings()),
	).build()
	if err != nil {
		return err
	}

	if d.config.Method == MethodLink {
		return runLink(list, force)
	}
	return runCopy(list, force)
}

func runLink(list []relation, force bool) error {
	links := make([]relation, 0, len(list))
	maybeAdd := func(v relation, err error) {
		if force {
			links = append(links, v)
		} else {
			pterm.Warning.Printfln("%s: %s", v.Dest.Path, err.Error())
		}
	}

	for _, v := range list { // create link list
		if v.Dest.NotExist {
			links = append(links, v)
			continue
		}
		if !v.Dest.isSymLink() { // if not symlink
			maybeAdd(v, fs.ErrExist)
			continue
		}
		if same, err := v.Dest.isSame(v.Src.Path); err != nil {
			return fmt.Errorf("%s: %w", v.Dest.Path, err)
		} else if !same {
			maybeAdd(v, fs.ErrExist)
		}
		// if src and dest are the same, do nothing
	}

	for _, v := range links { // do link
		// If the directory does not exist, create it
		dirPath := filepath.Dir(v.Dest.Path)
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return err
		}
		if !v.Dest.NotExist { // if file exists, remove it
			if err := os.Remove(v.Dest.Path); err != nil {
				return err
			}
		}
		if err := os.Symlink(v.Src.Path, v.Dest.Path); err != nil {
			return err
		}
		pterm.Success.Printfln("Symlink created. %s from %s", v.Dest.Path, v.Src.Path)
	}

	return nil
}

// If the file already exists, skip unless the force flag is true.
func runCopy(list []relation, force bool) error {
	copies := make([]relation, 0, len(list))
	maybeAdd := func(v relation, err error) {
		if force {
			copies = append(copies, v)
		} else {
			pterm.Warning.Printfln("%s: %s", v.Dest.Path, err.Error())
		}
	}

	for _, v := range list {
		if v.Dest.NotExist {
			copies = append(copies, v)
			continue
		}
		if v.Dest.isSymLink() {
			maybeAdd(v, errors.New("already linked"))
			continue
		}
		copies = append(copies, v)
	}

	for _, v := range copies {
		// If the directory does not exist, create it
		dirPath := filepath.Dir(v.Dest.Path)
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return err
		}

		if !v.Dest.NotExist { // if symlink exists, remove it
			if err := os.Remove(v.Dest.Path); err != nil {
				return err
			}
		}

		if err := copyFile(v.Src.Path, v.Dest.Path); err != nil {
			return err
		}

		pterm.Success.Printfln("File copied. %s from %s", v.Dest.Path, v.Src.Path)
	}

	return nil
}
