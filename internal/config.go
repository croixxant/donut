package internal

import (
	"errors"
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var config *viper.Viper

func GetDotDir() (string, error) {
	dotDir := config.GetString("dot_dir")
	if dotDir == "" {
		return "", errors.New("dor_dir is not defined")
	} else if f, err := os.Stat(dotDir); err != nil {
		return "", err
	} else if !f.IsDir() {
		return "", fmt.Errorf("%s is not directory", dotDir)
	}

	return dotDir, nil
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
