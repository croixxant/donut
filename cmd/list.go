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
		Args: cobra.NoArgs,
		RunE: List,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := internal.InitConfig(); err != nil {
				return err
			}
			cfg := internal.GetConfig()
			sourceDir, err := cfg.GetSrcDir()
			if err != nil {
				return err
			}
			if err := internal.InitMapConfig(sourceDir); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}

func List(cmd *cobra.Command, args []string) error {
	cfg := internal.GetConfig()
	sourceDir, err := cfg.GetSrcDir()
	if err != nil {
		return err
	}
	mapConfig := internal.GetMapConfig() // Get from config file
	destinationDir, err := mapConfig.GetDestDir()
	if err != nil {
		return err
	}
	list, err := newMap(sourceDir, destinationDir, mapConfig)
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
