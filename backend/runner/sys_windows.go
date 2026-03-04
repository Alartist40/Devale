//go:build windows
// +build windows

package runner

import (
	"os/exec"
	"syscall"
)

func setHideWindow(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}
