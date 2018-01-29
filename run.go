// +build ignore

package main

import (
	"flag"
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

var clean = flag.Bool("clean", false, "clean the previously generated files")
var resize = flag.Bool("resize", false, "resize content images to blog size")

func main() {
	flag.Parse()
	fmt.Println("building the blog...")
	if *clean {
		fmt.Println("cleaning up previous files")
		if err := runClean(); err != nil {
			log.Fatal(err)
		}
	}
	if *resize {
		fmt.Println("resizing images")
		if err := runResize(); err != nil {
			log.Fatal(err)
		}
		return
	}
	fmt.Println("runing hugo...")
	if err := runHugo(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("running go generate...")
	if err := runGoGenerate(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("running go build...")
	if err := runGoBuild(); err != nil {
		log.Fatal(err)
	}
}

func runClean() error {
	cmd, err := makeRmCmd("blog", "-rf", "public")
	if err != nil {
		errors.Wrap(err, "could not create rm cmd")
	}
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "%q did not run successfully", strings.Join(cmd.Args, " "))
	}
	cmd, err = makeRmCmd("blog", "static.go")
	if err != nil {
		errors.Wrap(err, "could not create rm cmd")
	}
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "%q did not run successfully", strings.Join(cmd.Args, " "))
	}
	return nil
}

type backupResizeImage struct {
	targetBasePath string
	backupBasePath string
	err            error
}

func (b *backupResizeImage) walkFn(path string, fi os.FileInfo, err error) error {
	if err != nil {
		return err
	} else if b.err != nil {
		return b.err
	}
	if fi.IsDir() {
		return nil
	}

	ext := strings.ToLower(filepath.Ext(path))
	if !strings.Contains(".jpg.png.gif", ext) {
		return nil
	}

	_, fileName := filepath.Split(path)
	matches, err := regexp.MatchString("^[1-9][0-9]*x[1-9][0-9]*_", fileName)
	if !matches {
		// file doesn't need to be resized so skip
		return nil
	} else if err != nil {
		return err
	}
	var maxWidth, maxHeight int
	_, err = fmt.Sscanf(fileName, "%dx%d_", &maxWidth, &maxHeight)

	b.backup(path)
	b.resize(path, maxWidth, maxHeight)
	if b.err != nil {
		return errors.Wrapf(err, "could not backup and resize image %q", path)
	}
	return nil
}

func (b *backupResizeImage) backup(path string) {
	if b.err != nil {
		return
	}
	rel, err := filepath.Rel(b.targetBasePath, path)
	if err != nil {
		b.err = err
		return
	}

	// find an unused filename for backup
	reldir, relname := filepath.Split(rel)
	ext := filepath.Ext(relname)
	name := strings.TrimSuffix(relname, ext)
	var backupPath string
	for count := 0; true; count++ {
		var suffix string
		if count > 0 {
			suffix = fmt.Sprintf("_%d", count)
		}
		backupPath = filepath.Join(b.backupBasePath, reldir, name+suffix+ext)
		_, err = os.Stat(backupPath)
		if err != nil {
			break
		}
	}

	dir, _ := filepath.Split(backupPath)
	if err := os.MkdirAll(dir, os.ModeDir|0775); err != nil {
		b.err = errors.Wrapf(err, "could not make backup dir %q", dir)
		return
	}
	inf, err := os.Open(path)
	if err != nil {
		b.err = errors.Wrapf(err, "could not open %q", path)
		return
	}
	defer inf.Close()
	outf, err := os.Create(backupPath)
	if err != nil {
		b.err = errors.Wrapf(err, "could not open file %q", backupPath)
		return
	}
	defer outf.Close()
	if _, err := io.Copy(outf, inf); err != nil {
		b.err = errors.Wrap(err, "could not create backup")
		return
	}
}

func (b *backupResizeImage) resize(path string, maxWidth, maxHeight int) {
	if b.err != nil {
		return
	}
	dimensions := fmt.Sprintf("%dx%d>", maxWidth, maxHeight)
	cmd := exec.Command("convert", path, "-resize", dimensions, "-quality", "85", "-interlace", "Line", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	b.err = cmd.Run()
}

func runResize() error {
	b := backupResizeImage{
		targetBasePath: "blog",
		backupBasePath: "backup",
	}
	if err := filepath.Walk(b.targetBasePath, b.walkFn); err != nil {
		return errors.Wrap(err, "could not walk blog/content")
	}
	return nil
}

func makeRmCmd(dir string, args ...string) (exec.Cmd, error) {
	path, err := exec.LookPath("rm")
	if err != nil {
		return exec.Cmd{}, errors.Wrapf(err, "coulld not find %q locally", "rm")
	}
	return exec.Cmd{
		Path:   path,
		Args:   append([]string{"rm"}, args...),
		Dir:    dir,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}, nil
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
		Stderr: os.Stderr,
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
