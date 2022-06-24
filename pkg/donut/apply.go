package donut

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/pterm/pterm"
)

func (d *Donut) Apply(force bool) error {
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

	if d.config.Method == MethodLink {
		return runLink(list, force)
	}
	return runCopy(list, force)
}

func runLink(list []relation, force bool) error {
	links := make([]relation, 0, len(list))
	maybeAdd := func(v relation, err error) {
		if force {
			links = append(links, v)
		} else {
			pterm.Warning.Printfln("%s: %s", v.Dest.Path, err.Error())
		}
	}

	for _, v := range list { // create link list
		if v.Dest.NotExist {
			links = append(links, v)
			continue
		}
		if !v.Dest.isSymLink() { // if not symlink
			maybeAdd(v, fs.ErrExist)
			continue
		}
		if same, err := v.Dest.isSame(v.Src.Path); err != nil {
			return fmt.Errorf("%s: %w", v.Dest.Path, err)
		} else if !same {
			maybeAdd(v, fs.ErrExist)
		}
		// if src and dest are the same, do nothing
	}

	for _, v := range links { // do link
		// If the directory does not exist, create it
		dirPath := filepath.Dir(v.Dest.Path)
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return err
		}
		if !v.Dest.NotExist { // if file exists, remove it
			if err := os.Remove(v.Dest.Path); err != nil {
				return err
			}
		}
		if err := os.Symlink(v.Src.Path, v.Dest.Path); err != nil {
			return err
		}
		pterm.Success.Printfln("Symlink created. %s from %s", v.Dest.Path, v.Src.Path)
	}

	return nil
}

// If the file already exists, skip unless the force flag is true.
func runCopy(list []relation, force bool) error {
	copies := make([]relation, 0, len(list))
	maybeAdd := func(v relation, err error) {
		if force {
			copies = append(copies, v)
		} else {
			pterm.Warning.Printfln("%s: %s", v.Dest.Path, err.Error())
		}
	}

	for _, v := range list {
		if v.Dest.NotExist {
			copies = append(copies, v)
			continue
		}
		if v.Dest.isSymLink() {
			maybeAdd(v, errors.New("already linked"))
			continue
		}
		copies = append(copies, v)
	}

	for _, v := range copies {
		// If the directory does not exist, create it
		dirPath := filepath.Dir(v.Dest.Path)
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return err
		}

		if !v.Dest.NotExist { // if symlink exists, remove it
			if err := os.Remove(v.Dest.Path); err != nil {
				return err
			}
		}

		if err := copyFile(v.Src.Path, v.Dest.Path); err != nil {
			return err
		}

		pterm.Success.Printfln("File copied. %s from %s", v.Dest.Path, v.Src.Path)
	}

	return nil
}
