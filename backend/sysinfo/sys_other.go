//go:build !windows
// +build !windows

package sysinfo

import (
	"os/exec"
)

func setHideWindow(cmd *exec.Cmd) {
	// No-op
}
