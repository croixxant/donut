package cmd

import (
	"github.com/spf13/cobra"

	"github.com/croixxant/donut/internal"
)

func InitConfig(_ *cobra.Command, _ []string) error {
	return internal.InitConfig(internal.WithFile(appName, cfgDirPaths...))
}

func InitConfigAndMapConfig(_ *cobra.Command, _ []string) error {
	if err := internal.InitConfig(internal.WithFile(appName, cfgDirPaths...)); err != nil {
		return err
	}
	cfg := internal.GetConfig()
	if err := internal.IsDir(cfg.SrcDir); err != nil {
		return err
	}
	if err := internal.InitMapConfig(internal.WithFile(mapConfigName, cfg.SrcDir)); err != nil {
		return err
	}
	return nil
}
