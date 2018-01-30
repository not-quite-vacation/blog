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
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "runs go generate in blog",
	Long:  `runs go generate in the blog directory`,
	Run:   runGoGenerate,
}

func init() {
	rootCmd.AddCommand(generateCmd)
}

func runGoGenerate(cmd *cobra.Command, args []string) {
	const gocmd = "go"
	path, err := exec.LookPath("go")
	if err != nil {
		log.Fatal(errors.Wrapf(err, "could not find %q locally", gocmd))
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
		log.Fatal(errors.Wrapf(err, "%q did not run successfully", strings.Join(gogen.Args, " ")))
	}
}
