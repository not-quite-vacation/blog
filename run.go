// +build ignore

package main

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

func main() {
	if err := runHugo(); err != nil {
		log.Fatal(err)
	}
	if err := runGoGenerate(); err != nil {
		log.Fatal(err)
	}
	if err := runGoBuild(); err != nil {
		log.Fatal(err)
	}
}

func runHugo() error {
	const hugocmd = "hugo"
	path, err := exec.LookPath(hugocmd)
	if err != nil {
		return errors.Wrapf(err, "could not find %q locally", hugocmd)
	}
	hugo := exec.Cmd{
		Path:   path,
		Dir:    "blog",
		Stdout: os.Stdout,
	}
	if err := hugo.Run(); err != nil {
		return errors.Wrapf(err, "%q did not run successfully", hugocmd)
	}
	return nil
}

func runGoGenerate() error {
	const gocmd = "go"
	path, err := exec.LookPath("go")
	if err != nil {
		return errors.Wrapf(err, "could not find %q locally", gocmd)
	}
	gogen := exec.Cmd{
		Path:   path,
		Args:   []string{gocmd, "generate"},
		Env:    append(os.Environ(), "GOOS=linux", "GOARCH=amd64"),
		Dir:    "blog",
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	if err := gogen.Run(); err != nil {
		return errors.Wrapf(err, "%q did not run successfully", strings.Join(gogen.Args, " "))
	}
	return nil
}

func runGoBuild() error {
	const gocmd = "go"
	path, err := exec.LookPath("go")
	if err != nil {
		return errors.Wrapf(err, "could not find %q locally", gocmd)
	}
	gobuild := exec.Cmd{
		Path:   path,
		Args:   []string{gocmd, "build", "-o", "serve"},
		Env:    os.Environ(),
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	if err := gobuild.Run(); err != nil {
		return errors.Wrapf(err, "%q did not run successfully", strings.Join(gobuild.Args, " "))
	}
	return nil
}
