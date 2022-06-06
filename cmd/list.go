package cmd

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/croixxant/donut/internal"
)

func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List files to be applied",
		Args:    cobra.NoArgs,
		PreRunE: InitConfigAndMapConfig,
		RunE: func(_ *cobra.Command, _ []string) error {
			return List()
		},
	}
	return cmd
}

func List() error {
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

	tableData := make([][]string, 0, len(list)+1) // add header capacity
	header := []string{"SOURCE", "DESTINATION"}
	if mapConfig.Method == internal.MethodLink {
		header = append([]string{"✔ "}, header...)
	}
	tableData = append(tableData, header)
	for _, v := range list {
		row := []string{v.Src.Path, v.Dest.Path}
		if mapConfig.Method == internal.MethodLink {
			var l string
			if !v.Dest.NotExist {
				if linked, err := v.Dest.IsSame(v.Src.Path); err != nil {
					return err
				} else if linked {
					l = "✔ "
				}
			}
			row = append([]string{l}, row...)
		}
		tableData = append(tableData, row)
	}

	if err := pterm.DefaultTable.
		WithHasHeader().
		WithData(tableData).
		WithBoxed().
		Render(); err != nil {
		return err
	}
	return nil
}
