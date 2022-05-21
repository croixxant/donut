package internal

import (
	"os"

	"github.com/mitchellh/mapstructure"
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
	Maps     []Map
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
	v := viper.New()

	v.SetConfigName(MapConfigName)
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
	return v.Unmarshal(&mapConfig.Data, viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
		ExpandEnvFunc(),
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
	)))
}

func (d *MapConfigData) AbsMaps(srcDir string, destDir string) []Map {
	s := make([]Map, 0, len(d.Maps))
	for _, v := range d.Maps {
		s = append(s, Map{
			Src:  Abs(v.Src, srcDir),
			Dest: Abs(v.Dest, destDir),
		})
	}
	return s
}
