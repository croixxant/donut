package cmd

import (
	"io"
	"os"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
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
