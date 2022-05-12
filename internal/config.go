package internal

import (
	"errors"
	"fmt"
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type (
	Config struct {
		Info  *ConfigData
		viper *viper.Viper
	}
	ConfigData struct {
		SrcDir string `mapstructure:"src_dir"`
	}
)

var config *Config = &Config{}

func GetConfig() *ConfigData {
	return config.Info
}

func InitConfig() error {
	v := viper.New()

	{ // initialize viper instance
		v.SetConfigName(AppName)
		v.AddConfigPath("$HOME")
		v.AddConfigPath("$XDG_CONFIG_HOME")
		if err := v.ReadInConfig(); err != nil {
			return err
		}
	}

	config.viper = v
	err := v.Unmarshal(&config.Info, viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
		ExpandEnvFunc(),
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
	)))

	return err
}

// GetSrcDir checks SrcDir is exists and is directory
func (d *ConfigData) GetSrcDir() (string, error) {
	if d.SrcDir == "" {
		return "", errors.New("src_dir is not defined")
	} else if fi, err := os.Stat(d.SrcDir); err != nil {
		return "", err
	} else if !fi.IsDir() {
		return "", fmt.Errorf("%s is not directory", d.SrcDir)
	}
	return d.SrcDir, nil
}
