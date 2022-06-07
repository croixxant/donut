package internal

import (
	"github.com/spf13/viper"
)

type (
	Config struct {
		Data  *ConfigData
		viper *viper.Viper
	}
	ConfigData struct {
		SrcDir   string `mapstructure:"src_dir"`
		DestDir  string `mapstructure:"dest_dir"`
		Excludes []string
		Method   Method
		Maps     map[string]string
	}
)

type Method string

const (
	MethodLink = "link"
	MethodCopy = "copy"
)

var config *Config = &Config{}

func GetConfig() *ConfigData {
	return config.Data
}

func WriteConfig(filename string) error {
	return config.viper.WriteConfigAs(filename)
}

func SetConfig(key string, value interface{}) (*ConfigData, error) {
	config.viper.Set(key, value)
	if err := config.viper.Unmarshal(&config.Data, viper.DecodeHook(defaultDecodeHookFunc)); err != nil {
		return nil, err
	}

	return config.Data, nil
}

func InitConfig(opts ...Option) error {
	v := viper.NewWithOptions(viper.KeyDelimiter("::"))

	for _, opt := range opts {
		if err := opt(v); err != nil {
			return err
		}
	}

	config.viper = v
	config.Data = &ConfigData{}
	return v.Unmarshal(&config.Data, viper.DecodeHook(defaultDecodeHookFunc))
}

func (d *ConfigData) AbsMaps(srcDir string, destDir string) map[string]string {
	m := make(map[string]string, len(d.Maps))
	for k, v := range d.Maps {
		m[Abs(k, srcDir)] = Abs(v, destDir)
	}
	return m
}
