package cmd

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/pterm/pterm"

	"github.com/croixxant/donut/internal"
)

func Copy(list []internal.Map) error {
	doList := make([]internal.Map, 0, len(list))
	for _, v := range list {
		if err := v.CanCopy(); err != nil {
			if errors.Is(err, internal.ErrAlreadyExists) {
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

		srcFile, err := os.Open(v.Src)
		if err != nil {
			return err
		}
		defer srcFile.Close()
		destFile, err := os.Create(v.Dest)
		if err != nil {
			return err
		}
		defer destFile.Close()
		b, err := io.ReadAll(srcFile)
		if err != nil {
			return err
		}
		_, err = destFile.Write(b)
		if err != nil {
			return err
		}

		pterm.Success.Printfln("File copied. %s from %s", v.Dest, v.Src)
	}

	return nil
}
