package cmd

import (
	"os"
	"path/filepath"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/croixxant/donut/internal"
)

func newInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [src_dir]",
		Short: "Generate the configuration file",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return Init(args[0])
		},
	}

	return cmd
}

func Init(srcDir string) error {
	if err := internal.InitConfig(internal.WithFile(appName, cfgDirPaths...)); err != nil {
		if err := internal.InitConfig(); err != nil {
			return err
		}
	}

	if _, err := internal.SetConfig("src_dir", srcDir); err != nil {
		return err
	}

	cfgPath := os.ExpandEnv(defaultConfigPath)
	dir := filepath.Dir(cfgPath)
	if err := internal.Mkdir(dir); err != nil {
		return err
	}

	if err := internal.WriteConfig(cfgPath); err != nil {
		return err
	}

	pterm.Success.Printfln("Configuration file created in %s", cfgPath)
	return nil
}
