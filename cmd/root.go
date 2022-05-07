package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/croixxant/donut/internal"
)

// rootCmd represents the base command when called without any subcommands
func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "donut",
		Version: internal.GetVersion(),
		Short:   "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
	examples and usage of using your application. For example:
	
	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra.Command, args []string) { },
		// SilenceErrors: true,
		SilenceUsage: true,
	}

	cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	cmd.AddCommand(newWhereCmd(), newListCmd())

	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := newRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
