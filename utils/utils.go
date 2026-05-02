package utils

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// universal utils that can be used by any package

// ------------------------------------
//
//	Check for the existence of an item in a slice
//
// ------------------------------------
func Contains[T comparable](slice []T, item T) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// ------------------------------------
//
//	Open a URL in the system's default browser in a non-blocking goroutine
//
// ------------------------------------
func OpenBrowser(url string) {
	go func() {
		var cmdExecutor *exec.Cmd

		switch runtime.GOOS {
		case "darwin":
			// macOS
			cmdExecutor = exec.Command("open", url)
		case "windows":
			// Windows
			cmdExecutor = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
		default:
			// Linux, BSD, WSL
			cmdExecutor = exec.Command("xdg-open", url)
		}

		cmdExecutor.Start()
	}()
}

// ------------------------------------
//
//	Return the platform-specific Downloads directory path, creating it if needed
//
// ------------------------------------
func GetDownloadsDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	var dir string
	switch runtime.GOOS {
	case "windows":
		dir = filepath.Join(home, "Downloads")
	case "darwin": // macOS
		dir = filepath.Join(home, "Downloads")
	case "linux":
		// Respect XDG_DOWNLOAD_DIR if set (most common on modern desktops)
		if d := os.Getenv("XDG_DOWNLOAD_DIR"); d != "" && filepath.IsAbs(d) {
			dir = d
		} else {
			dir = filepath.Join(home, "Downloads")
		}
	default:
		// Fallback for other Unix-like systems
		dir = filepath.Join(home, "Downloads")
	}

	// ensure the directory exists before returning
	makeDirErr := os.MkdirAll(dir, 0o755)
	if makeDirErr != nil {
		return "", makeDirErr
	}
	return dir, nil
}
