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
	"path/filepath"

	"github.com/spf13/cobra"
)

// renameCmd represents the rename command
var renameCmd = &cobra.Command{
	Use:   "rename",
	Short: "add a prefix to files",
	Long:  `rename all the files to add a prefix`,
	Run:   run,
}

var prefix *string

func init() {
	rootCmd.AddCommand(renameCmd)
	prefix = renameCmd.Flags().String("prefix", "", "prefix to add to the files")
}

func run(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		log.Fatal("you must supply files to rename")
	} else if *prefix == "" {
		log.Fatal("you must specify a prefix for the files")
	}

	for _, arg := range args {
		dir, file := filepath.Split(arg)
		if err := os.Rename(arg, filepath.Join(dir, *prefix+file)); err != nil {
			log.Fatal(err)
		}
	}
}
