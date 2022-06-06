package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/croixxant/donut/internal"
)

func newWhereCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "where",
		Short:   "Show dotfiles source directory",
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
