package donut

import (
	"io"

	"github.com/spf13/cobra"

	"github.com/croixxant/donut/pkg/donut"
)

func newWhereCmd(outStream, errStream io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "where",
		Short:   "Show dotfiles source directory",
		Args:    cobra.NoArgs,
		PreRunE: initConfig,
		RunE: func(cmd *cobra.Command, _ []string) error {
			d, err := donut.New(
				donut.WithConfig(donut.GetConfig()),
				donut.WithOut(outStream),
				donut.WithErr(errStream),
			)
			if err != nil {
				return err
			}
			return d.Where()
		},
	}

	cmd.SetOut(outStream)
	cmd.SetErr(errStream)

	return cmd
}
