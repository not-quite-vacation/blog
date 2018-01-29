// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// resizeCmd represents the resize command
var resizeCmd = &cobra.Command{
	Use:   "resize",
	Short: "resizes images based on filename",
	Long: `Walks through a directory finding images that need to be
resized. It will create a backup file of the original image in a specified
directory. Then resize the image.`,
	Run: runResize,
}

var baseDir, backupDir *string

func init() {
	rootCmd.AddCommand(resizeCmd)
	baseDir = resizeCmd.Flags().String("base-dir", "blog", "The directory to be search for images")
	backupDir = resizeCmd.Flags().String("backup-dir", "backup", "The directory to store the original images that are resized")
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

func runResize(cmd *cobra.Command, args []string) {
	b := backupResizeImage{
		targetBasePath: *baseDir,
		backupBasePath: *backupDir,
	}
	if err := filepath.Walk(b.targetBasePath, b.walkFn); err != nil {
		log.Fatal(errors.Wrapf(err, "could not walk %q", *baseDir))
	}
}
