//go:build !windows
// +build !windows

package runner

import (
	"os/exec"
)

func setHideWindow(cmd *exec.Cmd) {
	// No-op on non-windows
}
