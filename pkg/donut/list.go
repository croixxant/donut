package donut

import (
	"errors"

	"github.com/pterm/pterm"
)

func (d *Donut) List() error {
	if d.config == nil {
		return errors.New("no config")
	}

	list, err := newRelationsBuilder(
		d.config.SrcDir,
		d.config.DestDir,
		withExcludes(d.config.Excludes...),
		withMapping(d.config.absMappings()),
	).build()
	if err != nil {
		return err
	}

	tableData := make([][]string, 0, len(list)+1) // add header capacity
	header := []string{"SOURCE", "DESTINATION"}
	if d.config.Method == MethodLink {
		header = append([]string{"✔ "}, header...)
	}
	tableData = append(tableData, header)
	for _, v := range list {
		row := []string{v.Src.Path, v.Dest.Path}
		if d.config.Method == MethodLink {
			var l string
			if !v.Dest.NotExist {
				if linked, err := v.Dest.isSame(v.Src.Path); err != nil {
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
