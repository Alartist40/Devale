package runner

import (
	"os"
	"runtime"
)

func IsAdmin() bool {
	if runtime.GOOS != "windows" {
		return true // Assume true for dev on linux
	}

	f, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		return false
	}
	f.Close()
	return true
}
