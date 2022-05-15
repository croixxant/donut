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
	Files    []Map
}

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

func (d *MapConfigData) AbsFiles(srcDir string, destDir string) []Map {
	s := make([]Map, 0, len(d.Files))
	for _, v := range d.Files {
		s = append(s, Map{
			Src:  Abs(v.Src, srcDir),
			Dest: Abs(v.Dest, destDir),
		})
	}
	return s
}
