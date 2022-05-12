package internal

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type FileMapConfig struct {
	Excludes []string
	DestDir  string
	Files    []FileMap
}

type FileMap struct {
	Src  string
	Dest string
}

var fileMapConfig *viper.Viper

func GetFileMapConfig() (*FileMapConfig, error) {
	var cfg *FileMapConfig
	err := fileMapConfig.Unmarshal(&cfg)
	return cfg, err
}

func InitFileMapConfig(path string) error {
	fileMapConfig = viper.New()

	fileMapConfig.SetConfigName(FileMapConfigName)
	fileMapConfig.AddConfigPath(path)
	if err := fileMapConfig.ReadInConfig(); err != nil {
		return fmt.Errorf("Fatal error config file: %w \n", err)
	}
	fileMapConfig.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	fileMapConfig.WatchConfig()
	return nil
}
