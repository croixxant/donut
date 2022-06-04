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
		Use:   "init",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:
	
	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
		Args: cobra.ExactArgs(1),
		RunE: Init,
	}

	return cmd
}

var defaultConfigPath = "$HOME/.config/donut/donut.toml"

func Init(cmd *cobra.Command, args []string) error {
	if err := internal.InitConfig(internal.WithFile(appName, cfgDirPaths...)); err != nil {
		if err := internal.InitConfig(); err != nil {
			return err
		}
	}

	if _, err := internal.SetConfig("src_dir", args[0]); err != nil {
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
