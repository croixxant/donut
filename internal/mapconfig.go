package internal

import (
	"os"

	"github.com/spf13/viper"
)

type MapConfig struct {
	Data  *MapConfigData
	viper *viper.Viper
}

type MapConfigData struct {
	Excludes []string
	DestDir  string `mapstructure:"dest_dir"`
	Method   Method
	Maps     map[string]string
}

type Method string

const (
	MethodLink = "link"
	MethodCopy = "copy"
)

var mapConfig *MapConfig = &MapConfig{}

func GetMapConfig() *MapConfigData {
	return mapConfig.Data
}

func InitMapConfig(opts ...Option) error {
	v := viper.NewWithOptions(viper.KeyDelimiter("::"))

	{ // set default values
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		v.SetDefault("dest_dir", home)
		v.SetDefault("method", MethodCopy)
	}

	for _, opt := range opts {
		_ = opt(v)
	}

	mapConfig.viper = v
	return v.Unmarshal(&mapConfig.Data, viper.DecodeHook(defaultDecodeHookFunc))
}

func (d *MapConfigData) AbsMaps(srcDir string, destDir string) map[string]string {
	m := make(map[string]string, len(d.Maps))
	for k, v := range d.Maps {
		m[Abs(k, srcDir)] = Abs(v, destDir)
	}
	return m
}
