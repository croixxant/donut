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
		PreRunE: PreList,
	}
	return cmd
}

func PreList(cmd *cobra.Command, args []string) error {
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
	list := internal.NewMapBuilder(
		cfg.SrcDir, destDir, internal.WithExcludes(mapConfig.Excludes), internal.WithRemaps(remaps),
	).Build()

	tableData := make([][]string, 0, len(list)+1) // add header capacity
	tableData = append(tableData, []string{"SOURCE", "DESTINATION"})
	for _, v := range list {
		tableData = append(tableData, []string{v.Src, v.Dest})
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
