package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

func listInstalledFormula() (out []string, err error) {
	out, err = runAndGetStdout("brew", "list")
	if err != nil {
		return
	}
	return
}

func gen() error {
	if _, err := os.Stat(formulaYaml); err == nil {
		return fmt.Errorf("%s: already exists", formulaYaml)
	}

	tappedFormula, err := listTappedFormula()
	if err != nil {
		return err
	}

	installedFormula, err := listInstalledFormula()
	if err != nil {
		return err
	}

	fs := make([]Formula, len(installedFormula))
	for i, f := range installedFormula {
		fs[i].Name = f
	}
	bri := Brionac{
		Tap: tappedFormula,
		Brew: Brew{
			fs,
		},
	}
	c, err := yaml.Marshal(bri)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(formulaYaml, c, 0644)
	if err != nil {
		return err
	}

	return nil
}

func listTappedFormula() (out []string, err error) {
	out, err = runAndGetStdout("brew", "tap")
	if err != nil {
		return
	}
	return
}
