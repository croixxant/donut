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

func WriteConfig(filename string) error {
	return config.viper.WriteConfigAs(filename)
}

func SetConfig(key string, value interface{}) (*ConfigData, error) {
	config.viper.Set(key, value)
	if err := config.viper.Unmarshal(&config.Data, viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
		ExpandEnvFunc(),
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
	))); err != nil {
		return nil, err
	}

	return config.Data, nil
}

type Option func(v *viper.Viper) error

var CfgDirPaths = []string{"$XDG_CONFIG_HOME/" + AppName, "$HOME/.config/" + AppName, "$HOME"}

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
