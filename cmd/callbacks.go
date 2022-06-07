package cmd

import (
	"github.com/spf13/cobra"

	"github.com/croixxant/donut/internal"
)

func InitConfig(_ *cobra.Command, _ []string) error {
	return internal.InitConfig(internal.WithFile(appName, cfgDirPaths...))
}
