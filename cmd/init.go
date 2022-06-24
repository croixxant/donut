package cmd

import (
	"io"

	"github.com/spf13/cobra"

	"github.com/croixxant/donut/app"
)

func newInitCmd(outStream, errStream io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [src_dir]",
		Short: "Generate the configuration file",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, _ := cmd.Flags().GetString("config")
			d, _ := app.New(app.WithOut(outStream), app.WithErr(errStream))
			var srcDir string
			if len(args) > 0 {
				srcDir = args[0]
			}
			return d.Init(srcDir, cfgPath)
		},
	}

	cmd.SetOut(outStream)
	cmd.SetErr(errStream)

	return cmd
}
