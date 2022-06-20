package donut

import (
	"io"

	"github.com/spf13/cobra"

	"github.com/croixxant/donut/pkg/donut"
)

func newListCmd(outStream, errStream io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List files to be applied",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			cfgPath, _ := cmd.Flags().GetString("config")
			return donut.InitConfig(cfgPath)
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			d, err := donut.New(
				donut.WithConfig(donut.GetConfig()),
				donut.WithOut(outStream),
				donut.WithErr(errStream),
			)
			if err != nil {
				return err
			}
			return d.List()
		},
	}

	cmd.SetOut(outStream)
	cmd.SetErr(errStream)

	return cmd
}
