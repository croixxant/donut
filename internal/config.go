package internal

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var config *viper.Viper

func GetConfig() *viper.Viper {
	return config
}

func SetConfig() error {
	config = viper.New()

	config.SetConfigName(AppName)
	config.SetConfigType("toml")
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
