package main

import (
	"github.com/deckarep/golang-set"
	"reflect"
	"sort"
)

type Formulas struct {
	installed []string
	described []string
}

func clean() (err error) {
	var fs Formulas
	fs.installed, err = listInstalledFormula()
	if err != nil {
		return err
	}
	fs.described, err = listDescribedFormula()
	if err != nil {
		return err
	}

	if !fs.Equal() {
		if err = fs.Clean(); err != nil {
			return err
		}
	}

	return nil
}

func (f Formulas) Equal() bool {
	sort.Strings(f.installed)
	sort.Strings(f.described)
	return reflect.DeepEqual(f.installed, f.described)
}

func (f Formulas) Clean() (err error) {
	installedClasses := mapset.NewSet()
	for _, f := range f.installed {
		installedClasses.Add(f)
	}
	describedClasses := mapset.NewSet()
	for _, f := range f.described {
		describedClasses.Add(f)
	}

	for _, f := range installedClasses.Difference(describedClasses).ToSlice() {
		err = run([]string{"brew", "uninstall", f.(string)}, Blue)
	}

	if err != nil {
		return err
	}

	return nil
}
