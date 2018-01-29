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

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// hugoCmd represents the hugo command
var hugoCmd = &cobra.Command{
	Use:   "hugo",
	Short: "run hugo ",
	Long:  `runs hugo in the blog directory`,
	Run:   runHugo,
}

func init() {
	rootCmd.AddCommand(hugoCmd)
}

func runHugo(cmd *cobra.Command, args []string) {
	const (
		hugocmd = "hugo"
		dir     = "blog"
	)
	path, err := exec.LookPath(hugocmd)
	if err != nil {
		log.Fatal(errors.Wrapf(err, "could not find %q locally", hugocmd))
	}
	hugo := exec.Cmd{
		Path:   path,
		Dir:    dir,
		Args:   append([]string{"hugo"}, args...),
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	if err := hugo.Run(); err != nil {
		log.Fatal(errors.Wrapf(err, "%q did not run successfully", hugocmd))
	}
}
