package app

import (
	"io"
	"os"

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
