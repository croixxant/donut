package cmd

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"

	"github.com/croixxant/donut/internal"
)

type SyncMap struct {
	Src  string
	Dest string
}

var IgnoreFiles = []string{".git", ".gitignore"}

func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:
	
	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
		Args: cobra.NoArgs,
		RunE: List,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := internal.SetConfig(); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}

func List(cmd *cobra.Command, args []string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	dotDir, err := internal.GetDotDir()
	if err != nil {
		return err
	}
	entries, err := os.ReadDir(dotDir)
	if err != nil {
		return err
	}
	list, err := createSyncMap(entries, dotDir, home)
	if err != nil {
		return err
	}

	tableData := make([][]string, 0, len(list)+1)
	tableData = append(tableData, []string{"SOURCE", "DESTINATION"})
	for _, v := range list {
		tableData = append(tableData, []string{v.Src, v.Dest})
	}

	if err := pterm.DefaultTable.WithHasHeader().WithData(tableData).WithBoxed().Render(); err != nil {
		return err
	}
	return nil
}

func createSyncMap(entries []fs.DirEntry, srcDir, destDir string) (list []SyncMap, err error) {
	for _, v := range entries {
		name := v.Name()
		if slices.Contains(IgnoreFiles, name) {
			continue
		}
		srcPath := filepath.Join(srcDir, name)
		destPath := filepath.Join(destDir, name)
		if !v.IsDir() {
			list = append(list, SyncMap{
				Src:  srcPath,
				Dest: destPath,
			})
			continue
		}
		entries, err := os.ReadDir(srcPath)
		if err != nil {
			return []SyncMap{}, err
		}
		childList, err := createSyncMap(entries, srcPath, destPath)
		if err != nil {
			return []SyncMap{}, err
		}
		list = append(list, childList...)
	}
	return list, nil
}
