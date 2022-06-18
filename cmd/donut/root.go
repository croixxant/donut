package donut

import (
	"errors"
	"io"
	"os"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/croixxant/donut/pkg/donut"
)

// rootCmd represents the base command when called without any subcommands
func newRootCmd(outStream, errStream io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "donut",
		Version:      GetVersion(),
		Short:        "Tiny dotfiles management tool written in Go.",
		SilenceUsage: true,
	}

	cmd.SetOut(outStream)
	cmd.SetErr(errStream)

	cmd.PersistentFlags().StringP("config", "c", "", "location of config file")

	cmd.AddCommand(
		newInitCmd(outStream, errStream),
		newWhereCmd(outStream, errStream),
		newListCmd(outStream, errStream),
		newApplyCmd(outStream, errStream),
	)

	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := newRootCmd(os.Stdout, os.Stderr).Execute(); err != nil {
		pterm.Error.Println(err)
		os.Exit(1)
	}
}

// initConfig is a function that initialize the configuration.
// if -c option is set, use it as config file
// if -c option is not set, use default config file, and default config file is not exist, use default params
func initConfig(cmd *cobra.Command, _ []string) error {
	cfgPath, _ := cmd.Flags().GetString("config")
	if cfgPath != "" {
		_, err := donut.NewConfig(donut.WithFile(cfgPath))
		return err
	}
	if _, err := donut.NewConfig(donut.WithDefault(), donut.WithNameAndPath(donut.Name, donut.DefaultConfigDirs()...)); err != nil {
		if errors.As(err, &viper.ConfigFileNotFoundError{}) {
			return nil
		}
		return err
	}
	return nil
}
