package git

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gohyuhan/gitti/executor"
)

var (
	ErrIgnoreTracked      = errors.New("file is tracked; ignore has no effect")
	ErrIgnoreNotUntracked = errors.New("file is no longer untracked")
)

type IgnoreResult struct {
	IgnorePath     string
	AlreadyPresent bool
	Wrote          bool
	StillVisible   bool
}

// IgnorePattern anchors a literal path at its ignore file's directory.
func IgnorePattern(relativePath string) (string, error) {
	if relativePath == "" || relativePath == "." || filepath.IsAbs(relativePath) || strings.HasPrefix(relativePath, "/") || path.Clean(relativePath) != relativePath || relativePath == ".." || strings.HasPrefix(relativePath, "../") || strings.ContainsAny(relativePath, "\x00\r\n") || (filepath.Separator == '\\' && strings.Contains(relativePath, "\\")) {
		return "", fmt.Errorf("invalid relative file path")
	}
	var pattern strings.Builder
	pattern.WriteByte('/')
	for index := 0; index < len(relativePath); index++ {
		char := relativePath[index]
		if char == '*' || char == '?' || char == '[' || char == '\\' || (char == ' ' && index == len(relativePath)-1) {
			pattern.WriteByte('\\')
		}
		pattern.WriteByte(char)
	}
	return pattern.String(), nil
}

// IgnoreFile adds an untracked file to the nearest applicable .gitignore.
func (gf *GitFiles) IgnoreFile(relativePath string) (IgnoreResult, error) {
	var result IgnoreResult
	repoPath := gf.repoPath
	if repoPath == "" {
		return result, fmt.Errorf("repository path unavailable")
	}
	_, err := IgnorePattern(relativePath)
	if err != nil {
		return result, err
	}
	if !gf.gitProcessLock.CanProceedWithGitOps() {
		return result, errors.New(gf.gitProcessLock.OtherProcessRunningWarning())
	}
	defer gf.gitProcessLock.ReleaseGitOpsLock()

	trackedCmd := executor.RunGitCmdAt(repoPath, []string{"ls-files", "--cached", "-z", "--", ":(literal)" + relativePath}, false)
	tracked, err := trackedCmd.Output()
	if err != nil {
		return result, fmt.Errorf("git ls-files: %w", err)
	}
	if len(tracked) != 0 {
		return result, ErrIgnoreTracked
	}

	ignoreDir := path.Dir(relativePath)
	for ignoreDir != "." {
		candidate := filepath.Join(repoPath, filepath.FromSlash(ignoreDir), ".gitignore")
		info, err := os.Lstat(candidate)
		if err == nil && info.Mode().IsRegular() {
			break
		}
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return result, err
		}
		ignoreDir = path.Dir(ignoreDir)
	}
	result.IgnorePath = path.Join(ignoreDir, ".gitignore")
	ignorePath := filepath.Join(repoPath, filepath.FromSlash(result.IgnorePath))
	if info, err := os.Lstat(ignorePath); err == nil && !info.Mode().IsRegular() {
		return result, fmt.Errorf("%s is not a regular file", result.IgnorePath)
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return result, err
	}
	patternPath := relativePath
	if ignoreDir != "." {
		patternPath = strings.TrimPrefix(relativePath, ignoreDir+"/")
	}
	pattern, err := IgnorePattern(patternPath)
	if err != nil {
		return result, err
	}

	content, err := os.ReadFile(ignorePath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return result, err
	}
	alreadyPresent := false
	for _, line := range strings.Split(string(content), "\n") {
		if strings.TrimSuffix(line, "\r") == pattern {
			alreadyPresent = true
		}
	}
	status, err := ignoreFileStatus(repoPath, relativePath)
	if err != nil {
		return result, err
	}
	if len(status) == 0 {
		if alreadyPresent {
			result.AlreadyPresent = true
			return result, nil
		}
		return result, ErrIgnoreNotUntracked
	}
	if !strings.HasPrefix(string(status), "?? "+relativePath+"\x00") {
		return result, ErrIgnoreTracked
	}
	// A later negation in this file can cancel an earlier identical rule.
	// Append the rule again so it wins, rather than reporting an ineffective duplicate.
	file, err := os.OpenFile(ignorePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return result, err
	}
	newline := "\n"
	if bytes.Contains(content, []byte("\r\n")) {
		newline = "\r\n"
	}
	line := pattern + newline
	if len(content) > 0 && content[len(content)-1] != '\n' {
		line = newline + line
	}
	written, writeErr := file.WriteString(line)
	result.Wrote = written == len(line)
	closeErr := file.Close()
	if err := errors.Join(writeErr, closeErr); err != nil {
		return result, err
	}
	status, err = ignoreFileStatus(repoPath, relativePath)
	if err != nil {
		return result, err
	}
	if len(status) != 0 && !strings.HasPrefix(string(status), "?? "+relativePath+"\x00") {
		return result, ErrIgnoreTracked
	}
	result.StillVisible = len(status) != 0
	return result, nil
}

func ignoreFileStatus(repoPath, relativePath string) ([]byte, error) {
	cmd := executor.RunGitCmdAt(repoPath, []string{"status", "--porcelain", "-z", "--untracked-files=all", "--", ":(literal)" + relativePath}, false)
	status, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git status: %w", err)
	}
	return status, nil
}
