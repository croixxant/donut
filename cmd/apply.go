package cmd

import (
	"fmt"
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
		Args: cobra.NoArgs,
		RunE: Apply,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := internal.InitConfig(); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}

func Apply(cmd *cobra.Command, args []string) error {
	cfg, err := internal.GetConfig()
	if err != nil {
		return err
	}
	sourceDir, err := findSourceDir(cfg)
	if err != nil {
		return err
	}
	fileMapConfig, err := internal.GetFileMapConfig() // Get from config file
	if err != nil {
		return err
	}
	destinationDir, err := findDestinationDir(fileMapConfig)
	if err != nil {
		return err
	}
	list, err := newFileMaps(sourceDir, destinationDir, fileMapConfig.Excludes)
	if err != nil {
		return err
	}

	doList := make([]internal.FileMap, 0, len(list))
	for _, v := range list {
		f, err := os.Lstat(v.Dest)
		if err != nil {
			if os.IsNotExist(err) { // if Lstat() returns not exists error
				doList = append(doList, v)
				continue
			}
			return fmt.Errorf("%s: %w", v.Dest, err) // if Lstat() returns other error
		}

		if f.Mode()&os.ModeSymlink == 0 { // if not symlink
			pterm.Warning.Printfln("%s: %s", v.Dest, fs.ErrExist.Error())
			continue
		}

		link, err := os.Readlink(v.Dest)
		if err != nil { // if Readlink() returns error
			return fmt.Errorf("%s: %w", v.Dest, err)
		}
		if link != v.Src { // if link is not same as source path
			pterm.Warning.Printfln("%s: %s", v.Dest, fs.ErrExist.Error())
			continue
		}
	}

	for _, v := range doList {
		if err := os.Symlink(v.Src, v.Dest); err != nil {
			return err
		}
		pterm.Success.Printfln("Symlink created. %s from %s", v.Dest, v.Src)
	}

	return nil
}
