package main

import (
	"bufio"
	"errors"
	"github.com/daviddengcn/go-colortext"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

type Color int

const (
	None Color = Color(ct.None)
	Red  Color = Color(ct.Red)
	Blue Color = Color(ct.Blue)
)

func handleSignal() {
	sc := make(chan os.Signal, 10)
	signal.Notify(sc, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	go func() {
		<-sc
		ct.ResetColor()
		os.Exit(0)
	}()
}

var stdout = os.Stdout
var stderr = os.Stderr
var stdin = os.Stdin

func run(args []string, c Color) error {
	if len(args) == 0 {
		return errors.New("")
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Stdin = stdin
	ct.ChangeColor(ct.Color(c), true, ct.None, false)
	err := cmd.Run()
	ct.ResetColor()
	return err
}

func justRun(args []string) error {
	if len(args) == 0 {
		return errors.New("")
	}
	cmd := exec.Command(args[0], args[1:]...)
	return cmd.Run()
}

func runAndGetStdout(args ...string) (out []string, err error) {
	cmd := exec.Command(args[0], args[1:]...)
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		return
	}

	cmd.Start()

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		out = append(out, scanner.Text())
	}

	cmd.Wait()

	return
}

func getPath(args ...string) string {
	s, err := runAndGetStdout(args...)
	if err != nil {
		return ""
	}

	return filepath.Join(strings.Join(s, " "))
}
