package main

import (
	//"errors"
	"fmt"
	"github.com/cheggaaa/pb"
	"github.com/kyokomi/emoji"
	"github.com/mgutz/ansi"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	//"time"
)

type Brionac struct {
	Tap  []string
	Brew Brew
}

type Brew struct {
	Formula []Formula `yaml:"install"`
}

type Formula struct {
	Name string
	Args string
}

const formulaYaml = "formula.yaml"

func readFormulaYaml() (Brionac, error) {
	var bri Brionac

	if _, err := os.Stat(formulaYaml); err != nil {
		return bri, fmt.Errorf("%s: no such file or directory", formulaYaml)
	}

	buf, err := ioutil.ReadFile(formulaYaml)
	if err != nil {
		return bri, err
	}

	err = yaml.Unmarshal(buf, &bri)
	if err != nil {
		str := []byte(err.Error())
		assigned := regexp.MustCompile(`(line \d+)`)
		group := assigned.FindSubmatch(str)
		if len(group) != 0 {
			err = fmt.Errorf("Syntax Error at %s in %s", string(group[0]), formulaYaml)
		}
	}

	return bri, err
}

func install() error {
	d, err := readFormulaYaml()
	if err != nil {
		return err
	}

	defer func() {
		fmt.Fprintln(os.Stdout, emoji.Sprintf("\n:dragon: Brionac ended the %s.", state))
	}()

	for _, tap := range d.Tap {
		if alreadyTapped(tap) {
			continue
		}
		run([]string{"brew", "tap", tap}, Blue)
	}

	var (
		sc = 0
		fc = 0
	)
	var success map[string]interface{} = map[string]interface{}{
		"name":  []string{},
		"count": 0,
	}
	var failure map[string]interface{} = map[string]interface{}{
		"name":  []string{},
		"count": 0,
	}
	_ = success

	var installFormula []Formula
	for _, f := range d.Brew.Formula {
		if alreadyInstalled(f.Name) {
			sc++
		} else {
			installFormula = append(installFormula, f)
		}
	}

	if len(installFormula) == 0 {
		fmt.Fprintf(os.Stdout, "%d formulas are already installed.\n", len(d.Brew.Formula))
		return nil
	} else {
		fmt.Fprintf(os.Stdout, "%d formulas installing...\n", len(d.Brew.Formula))
	}

	if verbose {
		for _, brew := range installFormula {
			var args []string
			if brew.Args != "" {
				args = strings.Split(brew.Args, " ")
			}
			if err := brew.Install(args); err != nil {
				fc++
				failure["name"] = append(failure["name"].([]string), brew.Name)
			} else {
				sc++
			}
		}
	} else {
		bar := pb.StartNew(len(installFormula))
		bar.SetMaxWidth(80)
		for _, brew := range installFormula {
			var args []string
			if brew.Args != "" {
				args = strings.Split(brew.Args, " ")
			}
			bar.Increment()
			if err := brew.Install(args); err != nil {
				fc++
				failure["name"] = append(failure["name"].([]string), brew.Name)
			} else {
				sc++
			}
		}
		bar.FinishPrint("complete processing!\n")
	}

	fmt.Print(ansi.LightWhite)
	fmt.Fprintf(os.Stdout, "success: %d    failure: %d\n", sc, fc)
	fmt.Print(ansi.Reset)
	if len(failure["name"].([]string)) != 0 {
		fmt.Print(ansi.Red)
		fmt.Fprintf(os.Stdout, "errors: %q\n", failure["name"].([]string))
		fmt.Fprintf(os.Stdout, "  You should run `brew install <errored_fourmula>'\n")
		fmt.Print(ansi.Reset)
	}

	return nil
}

func alreadyTapped(name string) bool {
	tappedFormula, err := listTappedFormula()
	if err != nil {
		return true
	}
	for _, c := range tappedFormula {
		if c == name {
			return true
		}
	}
	return false
}

func alreadyInstalled(name string) bool {
	bin := filepath.Join(Cellar, name)
	if _, err := os.Stat(bin); err == nil {
		return true
	}
	return false
}

func (f Formula) Install(args []string) error {
	if alreadyInstalled(f.Name) {
		//return errors.New("already installed")
		return nil
	}

	cmdArgs := []string{"brew", "install"}
	cmdArgs = append(cmdArgs, args...)
	cmdArgs = append(cmdArgs, f.Name)

	//fmt.Printf("installing %s\n", f.Name)
	if verbose {
		return run(cmdArgs, Blue)
	}
	return justRun(cmdArgs)
}

func listDescribedFormula() ([]string, error) {
	var formulas []string

	d, err := readFormulaYaml()
	if err != nil {
		return formulas, err
	}

	for _, f := range d.Brew.Formula {
		formulas = append(formulas, f.Name)
	}

	return formulas, nil
}
