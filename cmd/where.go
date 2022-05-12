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
		Args: cobra.NoArgs,
		RunE: Where,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := internal.InitConfig(); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}

func Where(cmd *cobra.Command, _ []string) error {
	cfg, err := internal.GetConfig()
	if err != nil {
		return err
	}
	sourceDir, err := findSourceDir(cfg)
	if err != nil {
		return err
	}

	fmt.Println(sourceDir)
	return nil
}
