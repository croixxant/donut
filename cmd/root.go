package cmd

import (
	"os"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/croixxant/donut/internal"
)

// rootCmd represents the base command when called without any subcommands
func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "donut",
		Version:      internal.GetVersion(),
		Short:        "Tiny dotfiles management tool written in Go.",
		SilenceUsage: true,
	}

	cmd.AddCommand(newInitCmd(), newWhereCmd(), newListCmd(), newApplyCmd())

	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := newRootCmd().Execute(); err != nil {
		pterm.Error.Println(err)
		os.Exit(1)
	}
}
