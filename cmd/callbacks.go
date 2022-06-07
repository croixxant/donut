package cmd

import (
	"github.com/spf13/cobra"

	"github.com/croixxant/donut/internal"
)

func InitConfig(cmd *cobra.Command, _ []string) error {
	cfgPath, err := cmd.Flags().GetString("config")
	if err != nil {
		return err
	}
	if cfgPath == "" {
		return internal.InitConfig(internal.WithNameAndPath(appName, cfgDirPaths...))
	} else {
		return internal.InitConfig(internal.WithFile(cfgPath))
	}
}
