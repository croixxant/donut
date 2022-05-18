package internal

import (
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type (
	Config struct {
		Data  *ConfigData
		viper *viper.Viper
	}
	ConfigData struct {
		SrcDir string `mapstructure:"src_dir"`
	}
)

var config *Config = &Config{}

func GetConfig() *ConfigData {
	return config.Data
}

type Option func(v *viper.Viper) error

func WithFile(paths ...string) Option {
	return func(v *viper.Viper) error {
		for _, path := range paths {
			v.AddConfigPath(path)
		}
		if err := v.ReadInConfig(); err != nil {
			return err
		}
		return nil
	}
}

func InitConfig(opts ...Option) error {
	v := viper.New()

	v.SetConfigName(AppName)
	for _, opt := range opts {
		if err := opt(v); err != nil {
			return err
		}
	}

	config.viper = v
	return v.Unmarshal(&config.Data, viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
		ExpandEnvFunc(),
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
	)))
}
