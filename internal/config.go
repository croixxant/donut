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
	if err := config.viper.Unmarshal(&config.Data, viper.DecodeHook(defaultDecodeHookFunc)); err != nil {
		return nil, err
	}

	return config.Data, nil
}

func InitConfig(opts ...Option) error {
	v := viper.New()

	for _, opt := range opts {
		if err := opt(v); err != nil {
			return err
		}
	}

	config.viper = v
	return v.Unmarshal(&config.Data, viper.DecodeHook(defaultDecodeHookFunc))
}
