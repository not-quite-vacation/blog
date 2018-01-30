// +build windows

package cmd

import "os/exec"

func convertCommand(path, dimensions string) *exec.Cmd {
	return exec.Command("magick", "convert", path, "-resize", dimensions, "-quality", "85", "-interlace", "Line", path)
}
