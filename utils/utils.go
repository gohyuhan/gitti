package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"unicode"
	"unicode/utf8"
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

// the named escapes git uses when it C-quotes a path; any other control character is written in octal
var controlCharacterEscapes = map[rune]string{
	'\a': `\a`,
	'\b': `\b`,
	'\t': `\t`,
	'\n': `\n`,
	'\v': `\v`,
	'\f': `\f`,
	'\r': `\r`,
}

// ------------------------------------
//
//	Replace every control character (C0, DEL, C1) and every invalid UTF-8 byte in s with a
//	C escape the way git does (\n, \t, \033, \377), so s stays on one line and cannot send
//	escape sequences to the terminal. Printable characters, including non-ASCII, stay as they are.
//
// ------------------------------------
func EscapeControlCharacters(s string) string {
	if utf8.ValidString(s) && strings.IndexFunc(s, unicode.IsControl) < 0 {
		return s
	}

	var escaped strings.Builder
	for index := 0; index < len(s); {
		character, size := utf8.DecodeRuneInString(s[index:])
		isInvalidByte := character == utf8.RuneError && size == 1
		namedEscape, hasNamedEscape := controlCharacterEscapes[character]
		switch {
		case isInvalidByte || (unicode.IsControl(character) && !hasNamedEscape):
			for _, characterByte := range []byte(s[index : index+size]) {
				fmt.Fprintf(&escaped, `\%03o`, characterByte)
			}
		case hasNamedEscape:
			escaped.WriteString(namedEscape)
		default:
			escaped.WriteString(s[index : index+size])
		}
		index += size
	}
	return escaped.String()
}
