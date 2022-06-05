package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/croixxant/donut/internal"
)

func newWhereCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "where",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:
	
	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
		Args:    cobra.NoArgs,
		PreRunE: InitConfig,
		RunE: func(_ *cobra.Command, _ []string) error {
			return Where()
		},
	}
	return cmd
}

func Where() error {
	cfg := internal.GetConfig()
	if err := internal.IsDir(cfg.SrcDir); err != nil {
		return err
	}

	fmt.Println(cfg.SrcDir)
	return nil
}
