package api

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"github.com/fsnotify/fsnotify"

	"gitti/api/git"
	"gitti/settings"
)

type GitDaemon struct {
	RepoPath                       string
	Watcher                        *fsnotify.Watcher
	DebounceDur                    time.Duration
	GitFilesActiveRefreshDur       time.Duration
	GitFetchActiveRefreshDur       time.Duration
	isGitGeneralInfoPassiveRunning atomic.Bool
	isGitFilesActiveRunning        atomic.Bool
	isGitFetchActiveRunning        atomic.Bool
	WatcherTimer                   *time.Timer
	GitFilesActiveTimer            *time.Timer
	GitFetchActiveTimer            *time.Timer
	StopChannel                    chan struct{}
	ErrorLog                       []error
	UpdateChannel                  chan string // to communicate back to main thread for an update event
}

var GITDAEMON *GitDaemon

func InitGitDaemon(repoPath string, updateChannel chan string) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}

	debounce := time.Duration(settings.GITTICONFIGSETTINGS.FileWatcherDebounceMS) * time.Millisecond
	gitFilesActiveRefreshDur := time.Duration(settings.GITTICONFIGSETTINGS.GitFilesActiveRefreshDurationMS) * time.Millisecond
	gitFetchActiveRefreshDur := time.Duration(settings.GITTICONFIGSETTINGS.GitFetchDurationMS) * time.Millisecond
	gd := &GitDaemon{
		RepoPath:                 filepath.Join(repoPath, ".git"),
		Watcher:                  w,
		DebounceDur:              debounce,
		GitFilesActiveRefreshDur: gitFilesActiveRefreshDur,
		GitFetchActiveRefreshDur: gitFetchActiveRefreshDur,
		WatcherTimer:             time.NewTimer(debounce), // milliseconds
		GitFilesActiveTimer:      time.NewTimer(gitFilesActiveRefreshDur),
		GitFetchActiveTimer:      time.NewTimer(gitFetchActiveRefreshDur),
		StopChannel:              make(chan struct{}),
		ErrorLog:                 make([]error, 0),
		UpdateChannel:            updateChannel,
	}
	gd.isGitFilesActiveRunning.Store(false)
	gd.isGitFetchActiveRunning.Store(false)
	gd.isGitGeneralInfoPassiveRunning.Store(false)
	gd.WatcherTimer.Stop()
	gd.GitFilesActiveTimer.Stop()
	gd.GitFetchActiveTimer.Stop()
	gd.WatchPath()

	GITDAEMON = gd
}

func (gd *GitDaemon) WatchPath() {
	err := gd.Watcher.Add(gd.RepoPath)
	if err != nil {
		gd.ErrorLog = append(gd.ErrorLog, err)
	}
	err = filepath.WalkDir(gd.RepoPath, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			gd.Watcher.Add(path)
		}
		return nil
	})
	if err != nil {
		gd.ErrorLog = append(gd.ErrorLog, err)
	}

}

func (gd *GitDaemon) Start() {
	go func() {
		// Initial call to get info of git
		if gd.UpdateChannel != nil {
			GetUpdatedGitInfo(gd.UpdateChannel)
		}
		gd.GitFilesActiveTimer.Reset(gd.GitFilesActiveRefreshDur)
		gd.GitFetchActiveTimer.Reset(gd.GitFetchActiveRefreshDur)
		// loop to stay active
		for {
			select {
			case event := <-gd.Watcher.Events:
				if gd.isRelevantEvent(event) {
					gd.resetDebounce()
				}
			case err := <-gd.Watcher.Errors:
				fmt.Println("Watcher error:", err)

			case <-gd.WatcherTimer.C:
				go func() {
					if gd.isGitGeneralInfoPassiveRunning.CompareAndSwap(false, true) {
						defer gd.isGitGeneralInfoPassiveRunning.Store(false)
						GetUpdatedGitInfo(gd.UpdateChannel)
					}
				}()
			case <-gd.GitFilesActiveTimer.C:
				// reset first to avoid losing ticks, then run work in background
				gd.GitFilesActiveTimer.Reset(gd.GitFilesActiveRefreshDur)
				go func() {
					if gd.isGitFilesActiveRunning.CompareAndSwap(false, true) {
						// Mark as running
						defer gd.isGitFilesActiveRunning.Store(false)

						git.GITFILES.GetGitFilesStatus()
						gd.UpdateChannel <- git.GIT_FILES_STATUS_UPDATE
					}
				}()
			case <-gd.GitFetchActiveTimer.C:
				// reset immediately; git fetch operation TBD
				gd.GitFetchActiveTimer.Reset(gd.GitFetchActiveRefreshDur)
				// go func() {

				// }()
			case <-gd.StopChannel:
				gd.Watcher.Close()
				return
			}
		}
	}()
}

func (gd *GitDaemon) resetDebounce() {
	if !gd.WatcherTimer.Stop() {
		select {
		case <-gd.WatcherTimer.C:
		default:
		}
	}
	gd.WatcherTimer.Reset(gd.DebounceDur)
}

func (gd *GitDaemon) isRelevantEvent(event fsnotify.Event) bool {
	// Only watch .git subpaths
	if !strings.Contains(event.Name, filepath.Join(gd.RepoPath)) {
		return false
	}

	// Ignore lock and temp files that git touches rapidly
	base := filepath.Base(event.Name)
	if strings.HasSuffix(base, ".lock") || base == "FETCH_HEAD" {
		return false
	}

	// Handle new directories
	if event.Op&fsnotify.Create == fsnotify.Create {
		fi, err := os.Stat(event.Name)
		if err == nil && fi.IsDir() {
			filepath.WalkDir(event.Name, func(path string, d fs.DirEntry, err error) error {
				if err == nil && d.IsDir() {
					_ = gd.Watcher.Add(path)
				}
				return nil
			})
		}
		return true
	}

	// Trigger only for relevant ops
	if event.Op&(fsnotify.Write|fsnotify.Remove|fsnotify.Rename) != 0 {
		return true
	}

	return false
}

func (gd *GitDaemon) Stop() {
	close(gd.StopChannel)
}
