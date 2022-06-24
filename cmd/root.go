package cmd

import (
	"io"
	"os"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/croixxant/donut/app"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := newRootCmd(os.Stdout, os.Stderr).Execute(); err != nil {
		pterm.Error.Println(err)
		os.Exit(1)
	}
}

// rootCmd represents the base command when called without any subcommands
func newRootCmd(outStream, errStream io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "donut",
		Version:      GetVersion(),
		Short:        "Tiny dotfiles management tool written in Go.",
		SilenceUsage: true,
	}

	cmd.SetOut(outStream)
	cmd.SetErr(errStream)

	cmd.PersistentFlags().StringP("config", "c", "", "location of config file")

	cmd.AddCommand(
		newInitCmd(outStream, errStream),
		newWhereCmd(outStream, errStream),
		newListCmd(outStream, errStream),
		newApplyCmd(outStream, errStream),
	)

	return cmd
}

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

func newWhereCmd(outStream, errStream io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "where",
		Short: "Show dotfiles source directory",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			cfgPath, _ := cmd.Flags().GetString("config")
			return app.InitConfig(cfgPath)
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			d, err := app.New(
				app.WithConfig(app.GetConfig()),
				app.WithOut(outStream),
				app.WithErr(errStream),
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

func newListCmd(outStream, errStream io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List files to be applied",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			cfgPath, _ := cmd.Flags().GetString("config")
			return app.InitConfig(cfgPath)
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			d, err := app.New(
				app.WithConfig(app.GetConfig()),
				app.WithOut(outStream),
				app.WithErr(errStream),
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

func newApplyCmd(outStream, errStream io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply files from source to destination",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			cfgPath, _ := cmd.Flags().GetString("config")
			return app.InitConfig(cfgPath)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			force, _ := cmd.Flags().GetBool("force")
			d, err := app.New(
				app.WithConfig(app.GetConfig()),
				app.WithOut(outStream),
				app.WithErr(errStream),
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
