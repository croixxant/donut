package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

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
		PreRunE: InitConfigAndMapConfig,
		RunE: func(cmd *cobra.Command, args []string) error {
			force, _ := cmd.Flags().GetBool("force")
			return Apply(force)
		},
	}

	cmd.Flags().BoolP("force", "f", false, "Help message for toggle")

	return cmd
}

func Apply(force bool) error {
	cfg := internal.GetConfig()
	if err := internal.IsDir(cfg.SrcDir); err != nil {
		return err
	}
	mapConfig := internal.GetMapConfig() // Get from config file
	destDir, err := internal.DirOrHome(mapConfig.DestDir)
	if err != nil {
		return err
	}

	remaps := mapConfig.AbsMaps(cfg.SrcDir, destDir)
	excludes := append(ignores, mapConfig.Excludes...)
	list, err := internal.NewMapBuilder(
		cfg.SrcDir, destDir, internal.WithExcludes(excludes...), internal.WithRemaps(remaps),
	).Build()
	if err != nil {
		return err
	}

	if mapConfig.Method == internal.MethodLink {
		return Link(list, force)
	}
	return Copy(list, force)
}

func Link(list []internal.Map, force bool) error {
	links := make([]internal.Map, 0, len(list))
	maybeAdd := func(v internal.Map, err error) {
		if force {
			links = append(links, v)
		} else {
			pterm.Warning.Printfln("%s: %s", v.Dest.Path, err.Error())
		}
	}

	for _, v := range list { // create link list
		if v.Dest.NotExist {
			links = append(links, v)
			continue
		}
		if !v.Dest.IsSymLink() { // if not symlink
			maybeAdd(v, fs.ErrExist)
			continue
		}
		if same, err := v.Dest.IsSame(v.Src.Path); err != nil {
			return fmt.Errorf("%s: %w", v.Dest.Path, err)
		} else if !same {
			maybeAdd(v, fs.ErrExist)
		}
		// if src and dest are the same, do nothing
	}

	for _, v := range links { // do link
		// If the directory does not exist, create it
		dirPath := filepath.Dir(v.Dest.Path)
		if err := internal.Mkdir(dirPath); err != nil {
			return err
		}
		if !v.Dest.NotExist { // if file exists, remove it
			if err := os.Remove(v.Dest.Path); err != nil {
				return err
			}
		}
		if err := os.Symlink(v.Src.Path, v.Dest.Path); err != nil {
			return err
		}
		pterm.Success.Printfln("Symlink created. %s from %s", v.Dest.Path, v.Src.Path)
	}

	return nil
}

// If the file already exists, skip unless the force flag is true.
func Copy(list []internal.Map, force bool) error {
	copies := make([]internal.Map, 0, len(list))
	maybeAdd := func(v internal.Map, err error) {
		if force {
			copies = append(copies, v)
		} else {
			pterm.Warning.Printfln("%s: %s", v.Dest.Path, err.Error())
		}
	}

	for _, v := range list {
		if v.Dest.NotExist {
			copies = append(copies, v)
			continue
		}
		if v.Dest.IsSymLink() {
			maybeAdd(v, errors.New("already linked"))
			continue
		}
		copies = append(copies, v)
	}

	for _, v := range copies {
		// If the directory does not exist, create it
		dirPath := filepath.Dir(v.Dest.Path)
		if err := internal.Mkdir(dirPath); err != nil {
			return err
		}

		if !v.Dest.NotExist { // if symlink exists, remove it
			if err := os.Remove(v.Dest.Path); err != nil {
				return err
			}
		}

		if err := internal.Copy(v.Src.Path, v.Dest.Path); err != nil {
			return err
		}

		pterm.Success.Printfln("File copied. %s from %s", v.Dest.Path, v.Src.Path)
	}

	return nil
}
