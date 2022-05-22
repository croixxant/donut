package cmd

import (
	"errors"
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
		RunE:    Apply,
		PreRunE: PreApply,
	}

	cmd.Flags().BoolP("force", "f", false, "Help message for toggle")

	return cmd
}

func PreApply(cmd *cobra.Command, args []string) error {
	if err := internal.InitConfig(internal.WithFile(internal.CfgDirPaths...)); err != nil {
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

	remaps := mapConfig.AbsMaps(cfg.SrcDir, destDir)
	list := internal.NewMapBuilder(
		cfg.SrcDir, destDir, internal.WithExcludes(mapConfig.Excludes), internal.WithRemaps(remaps),
	).Build()

	force, _ := cmd.Flags().GetBool("force")
	if mapConfig.Method == internal.MethodLink {
		return Link(list, force)
	}
	return Copy(list, force)
}

func Link(list []internal.Map, force bool) error {
	links := make([]internal.Map, 0, len(list))
	removes := make(map[string]bool, len(list))
	for _, v := range list {
		if err := v.CanLink(); err != nil {
			if errors.Is(err, internal.ErrAlreadyLinked) {
				continue
			} else if errors.Is(err, fs.ErrExist) {
				if force {
					links = append(links, v)
					removes[v.Dest] = true
				} else {
					pterm.Warning.Println(err)
				}
				continue
			}
			return err
		}
		links = append(links, v)
		removes[v.Dest] = false
	}

	for _, v := range links {
		// If the directory does not exist, create it
		dirPath := filepath.Dir(v.Dest)
		if err := internal.Mkdir(dirPath); err != nil {
			return err
		}

		if removes[v.Dest] { // if file exists, remove it
			if err := os.Remove(v.Dest); err != nil {
				return err
			}
		}

		if err := os.Symlink(v.Src, v.Dest); err != nil {
			return err
		}
		pterm.Success.Printfln("Symlink created. %s from %s", v.Dest, v.Src)
	}

	return nil
}

// If the file already exists, skip unless the force flag is true.
func Copy(list []internal.Map, force bool) error {
	copies := make([]internal.Map, 0, len(list))
	removes := make(map[string]bool, len(list))

	for _, v := range list {
		if err := v.CanCopy(); err != nil {
			if errors.Is(err, internal.ErrAlreadyLinked) {
				if force {
					copies = append(copies, v)
					removes[v.Dest] = true
				} else {
					pterm.Warning.Println(err)
				}
				continue
			}
			return err
		}
		copies = append(copies, v)
		removes[v.Dest] = false
	}

	for _, v := range copies {
		// If the directory does not exist, create it
		dirPath := filepath.Dir(v.Dest)
		if err := internal.Mkdir(dirPath); err != nil {
			return err
		}

		if removes[v.Dest] { // if symlink exists, remove it
			if err := os.Remove(v.Dest); err != nil {
				return err
			}
		}

		if err := internal.Copy(v.Src, v.Dest); err != nil {
			return err
		}

		pterm.Success.Printfln("File copied. %s from %s", v.Dest, v.Src)
	}

	return nil
}
