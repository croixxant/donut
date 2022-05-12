package internal

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	SrcDir string `mapstructure:"src_dir"`
}

var config *viper.Viper

func GetConfig() (*Config, error) {
	var cfg *Config
	err := config.Unmarshal(&cfg)
	return cfg, err
}

func InitConfig() error {
	config = viper.New()

	config.SetConfigName(AppName)
	config.AddConfigPath("$HOME")
	config.AddConfigPath("$XDG_CONFIG_HOME")
	config.AddConfigPath(".")
	if err := config.ReadInConfig(); err != nil {
		return fmt.Errorf("Fatal error config file: %w \n", err)
	}
	config.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	config.WatchConfig()
	return nil
}
