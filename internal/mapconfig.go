package internal

import (
	"os"
	"path/filepath"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type MapConfig struct {
	Info  *MapConfigData
	viper *viper.Viper
}

type MapConfigData struct {
	Excludes []string
	DestDir  string `mapstructure:"dest_dir"`
	Files    []Map
}

type Map struct {
	Src  string
	Dest string
}

var mapConfig *MapConfig = &MapConfig{}

func GetMapConfig() *MapConfigData {
	return mapConfig.Info
}

func InitMapConfig(path string) error {
	v := viper.New()

	{ // initialize viper instance
		v.SetConfigName(MapConfigName)
		v.AddConfigPath(path)
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		v.SetDefault("dest_dir", home)
		_ = v.ReadInConfig()
	}

	mapConfig.viper = v

	err := v.Unmarshal(&mapConfig.Info, viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
		ExpandEnvFunc(),
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
	)))

	return err
}

// GetDestDir checks DestDir is exists and is directory.
// if DestDir is not exist or is not directory, then returns home directory
func (d *MapConfigData) GetDestDir() (string, error) {
	if d.DestDir == "" {
		return os.UserHomeDir()
	} else if fi, err := os.Stat(d.DestDir); err != nil {
		return os.UserHomeDir()
	} else if !fi.IsDir() {
		return os.UserHomeDir()
	}
	return d.DestDir, nil
}

func (d *MapConfigData) Map(srcDir string) map[string]string {
	m := map[string]string{}
	for _, v := range d.Files {
		m[joinAndClean(v.Src, srcDir)] = joinAndClean(v.Dest, d.DestDir)
	}
	return m
}

func joinAndClean(path, baseDir string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Clean(filepath.Join(baseDir, path))
}
