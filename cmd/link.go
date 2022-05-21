package cmd

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/pterm/pterm"

	"github.com/croixxant/donut/internal"
)

func Link(list []internal.Map) error {
	doList := make([]internal.Map, 0, len(list))
	for _, v := range list {
		if err := v.CanLink(); err != nil {
			if errors.Is(err, internal.ErrAlreadyLinked) {
				continue
			} else if errors.Is(err, fs.ErrExist) {
				pterm.Warning.Println(err)
				continue
			}
			return err
		}
		doList = append(doList, v)
	}

	for _, v := range doList {
		dirPath := filepath.Dir(v.Dest)
		if err := internal.IsDir(dirPath); err != nil {
			if !os.IsNotExist(err) {
				return err
			}
			if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
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
