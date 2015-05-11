package main

import (
	"fmt"
	"github.com/mgutz/ansi"
	"os"
	"os/exec"
)

func initBrionac() error {
	fmt.Print(ansi.LightWhite)
	fmt.Fprintln(os.Stdout, "Initializing Brionac...")
	fmt.Print(ansi.Reset)

	homebrewInstaller := `ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"`
	homebrewPath, err := exec.LookPath("brew")
	if err != nil {
		return fmt.Errorf("command not found -> brew")
	}

	fmt.Fprintf(os.Stdout, "  not found brew, %s\n", homebrewInstaller)
	fmt.Fprintf(os.Stdout, "  not found brew, %s\n", homebrewPath)

	return nil
}
