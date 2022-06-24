package donut

import (
	"errors"

	"github.com/spf13/viper"
)

var viperInstance = viper.New()

func InitConfig(cfgPath string) error {
	if cfgPath != "" {
		_, err := NewConfig(WithFile(cfgPath))
		return err
	}
	if _, err := NewConfig(WithDefault(), WithNameAndPath(Name, DefaultConfigDirs()...)); err != nil {
		if errors.As(err, &viper.ConfigFileNotFoundError{}) {
			return nil
		}
		return err
	}
	return nil
}

func NewConfig(opts ...ConfigOption) (*viper.Viper, error) {
	viperInstance = viper.NewWithOptions(viper.KeyDelimiter("::"))

	for _, opt := range opts {
		if err := opt(viperInstance); err != nil {
			return nil, err
		}
	}

	return viperInstance, nil
}

func GetConfig() *viper.Viper {
	return viperInstance
}

func WriteConfig(filename string) error {
	return viperInstance.WriteConfigAs(filename)
}

type Config struct {
	SrcDir   string            `mapstructure:"src_dir"`
	DestDir  string            `mapstructure:"dest_dir"`
	Excludes []string          `mapstructure:"excludes"`
	Method   Method            `mapstructure:"method"`
	Mappings map[string]string `mapstructure:"mappings"`
}

type Method string

const (
	MethodLink = "link"
	MethodCopy = "copy"
)

func (d *Config) absMappings() map[string]string {
	m := make(map[string]string, len(d.Mappings))
	for k, v := range d.Mappings {
		m[absPath(k, d.SrcDir)] = absPath(v, d.DestDir)
	}
	return m
}

type ConfigOption func(v *viper.Viper) error

func WithFile(file string) ConfigOption {
	return func(v *viper.Viper) error {
		v.SetConfigFile(file)
		if err := v.ReadInConfig(); err != nil {
			return err
		}
		return nil
	}
}

func WithNameAndPath(name string, paths ...string) ConfigOption {
	return func(v *viper.Viper) error {
		v.SetConfigName(name)
		for _, path := range paths {
			v.AddConfigPath(path)
		}
		if err := v.ReadInConfig(); err != nil {
			return err
		}
		return nil
	}
}

func WithData(data map[string]interface{}) ConfigOption {
	return func(v *viper.Viper) error {
		for k, d := range data {
			v.Set(k, d)
		}
		return nil
	}
}

func WithDefault() ConfigOption {
	return func(v *viper.Viper) error {
		v.SetDefault("src_dir", DefaultSourceDir())
		v.SetDefault("dest_dir", UserHomeDir)
		v.SetDefault("method", MethodCopy)
		return nil
	}
}
