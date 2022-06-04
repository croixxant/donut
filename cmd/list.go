package cmd

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/croixxant/donut/internal"
)

func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:
	
	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
		Args:    cobra.NoArgs,
		RunE:    List,
		PreRunE: InitConfigAndMapConfig,
	}
	return cmd
}

func List(cmd *cobra.Command, args []string) error {
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
		cfg.SrcDir, destDir, internal.WithExcludes(excludes), internal.WithRemaps(remaps),
	).Build()
	if err != nil {
		return err
	}

	tableData := make([][]string, 0, len(list)+1) // add header capacity
	tableData = append(tableData, []string{"SOURCE", "DESTINATION"})
	for _, v := range list {
		tableData = append(tableData, []string{v.Src.Path, v.Dest.Path})
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
