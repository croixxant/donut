package donut

import (
	"io"

	"github.com/spf13/cobra"

	"github.com/croixxant/donut/pkg/donut"
)

func newApplyCmd(outStream, errStream io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply files from source to destination",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			cfgPath, _ := cmd.Flags().GetString("config")
			return donut.InitConfig(cfgPath)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			force, _ := cmd.Flags().GetBool("force")
			d, err := donut.New(
				donut.WithConfig(donut.GetConfig()),
				donut.WithOut(outStream),
				donut.WithErr(errStream),
			)
			if err != nil {
				return err
			}
			return d.Apply(force)
		},
	}

	cmd.SetOut(outStream)
	cmd.SetErr(errStream)

	cmd.Flags().BoolP("force", "f", false, "Force the application")

	return cmd
}
