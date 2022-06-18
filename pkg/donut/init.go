package donut

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pterm/pterm"
	"github.com/spf13/viper"
)

func (d *Donut) Init(srcDir, cfgPath string) (err error) {
	var v *viper.Viper
	if cfgPath == "" {
		_, err = NewConfig(WithDefault(), WithNameAndPath(Name, DefaultConfigDirs()...))
		cfgPath = DefaultConfigFile()
	} else {
		_, err = NewConfig(WithFile(cfgPath))
	}
	if err == nil {
		return errors.New("no need to init")
	}
	if !errors.As(err, &viper.ConfigFileNotFoundError{}) && !os.IsNotExist(err) {
		return fmt.Errorf("config file already exists, but error: %w", err)
	}

	v, _ = NewConfig(WithDefault())

	if srcDir != "" {
		v.Set("src_dir", srcDir)
	}
	if err := os.MkdirAll(filepath.Dir(cfgPath), os.ModePerm); err != nil {
		return err
	}
	if err := WriteConfig(cfgPath); err != nil {
		return err
	}

	pterm.Success.Printfln("Configuration file created in %s", cfgPath)
	return nil
}
