package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/croixxant/donut/internal"
	"github.com/spf13/cobra"
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
			if err := internal.SetConfig(); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}

func Where(cmd *cobra.Command, _ []string) error {
	dotDir := internal.GetConfig().GetString("dot_dir")
	if dotDir == "" {
		return errors.New("dor_dir is not defined")
	} else if f, err := os.Stat(dotDir); err != nil {
		return err
	} else if !f.IsDir() {
		return fmt.Errorf("%s is not directory", dotDir)
	}

	fmt.Println(dotDir)
	return nil
}