package donut

import (
	"errors"
	"fmt"
)

func (d *Donut) Where() error {
	if d.config == nil {
		return errors.New("no config")
	}
	fmt.Fprint(d.out, d.config.SrcDir)
	return nil
}
