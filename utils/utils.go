package utils

import (
	"os/exec"
	"runtime"
)

// universal utils that can be used by any package

// Contains is a generic helper function to check for the existence of an item in a slice.
func Contains[T comparable](slice []T, item T) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func OpenBrowser(url string) {
	go func() {
		var cmd *exec.Cmd

		switch runtime.GOOS {
		case "darwin":
			// macOS
			cmd = exec.Command("open", url)
		case "windows":
			// Windows
			cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
		default:
			// Linux, BSD, WSL
			cmd = exec.Command("xdg-open", url)
		}

		cmd.Start()
	}()
}
