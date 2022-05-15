package cmd

import (
	"errors"
	"io/fs"
	"os"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/croixxant/donut/internal"
)

func newApplyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apply",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:
	
	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
		Args:    cobra.NoArgs,
		RunE:    Apply,
		PreRunE: PreApply,
	}
	return cmd
}

func PreApply(cmd *cobra.Command, args []string) error {
	if err := internal.InitConfig(internal.WithFile("$HOME", "$XDG_CONFIG_HOME")); err != nil {
		return err
	}
	cfg := internal.GetConfig()
	if err := internal.IsDir(cfg.SrcDir); err != nil {
		return err
	}
	if err := internal.InitMapConfig(internal.WithFile(cfg.SrcDir)); err != nil {
		return err
	}
	return nil
}

func Apply(cmd *cobra.Command, args []string) error {
	cfg := internal.GetConfig()
	if err := internal.IsDir(cfg.SrcDir); err != nil {
		return err
	}
	mapConfig := internal.GetMapConfig() // Get from config file
	destDir, err := internal.DirOrHome(mapConfig.DestDir)
	if err != nil {
		return err
	}

	remaps := mapConfig.AbsFiles(cfg.SrcDir, destDir)
	list := internal.NewMapBuilder(
		cfg.SrcDir, destDir, internal.WithExcludes(mapConfig.Excludes), internal.WithRemaps(remaps),
	).Build()

	doList := make([]internal.Map, 0, len(list))
	for _, v := range list {
		if err := v.CanLink(); err != nil {
			if errors.Is(err, internal.ErrAlreadyLinked) {
				continue
			} else if errors.Is(err, fs.ErrExist) {
				pterm.Warning.Println(err)
				continue
			}
			return err
		}
		doList = append(doList, v)
	}

	for _, v := range doList {
		if err := os.Symlink(v.Src, v.Dest); err != nil {
			return err
		}
		pterm.Success.Printfln("Symlink created. %s from %s", v.Dest, v.Src)
	}

	return nil
}
