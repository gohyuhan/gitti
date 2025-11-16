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
	repoPath                            string
	watcher                             *fsnotify.Watcher
	debounceDur                         time.Duration
	gitFilesActiveRefreshDur            time.Duration
	gitRemoteSyncStatusActiveRefreshDur time.Duration
	isGitBranchPassiveRunning           atomic.Bool
	isGitFilesPassiveActiveRunning      atomic.Bool
	isGitStashPassiveRunning            atomic.Bool
	isGitRemoteSyncStatusActiveRunning  atomic.Bool
	watcherTimer                        *time.Timer
	gitFilesActiveTimer                 *time.Timer
	gitRemoteSyncStatusActiveTimer      *time.Timer
	stopChannel                         chan struct{}
	errorLog                            []error
	updateChannel                       chan string // to communicate back to main thread for an update event
	gitOperations                       *GitOperations
}

var GITDAEMON *GitDaemon

func InitGitDaemon(absoluteGitPath string, updateChannel chan string, gitOperations *GitOperations) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}

	debounce := time.Duration(settings.GITTICONFIGSETTINGS.FileWatcherDebounceMS) * time.Millisecond
	gitFilesActiveRefreshDur := time.Duration(settings.GITTICONFIGSETTINGS.GitFilesActiveRefreshDurationMS) * time.Millisecond
	gitRemoteSyncStatusActiveRefreshDur := time.Duration(settings.GITTICONFIGSETTINGS.GitRemoteSyncStatusDurationMS) * time.Millisecond
	gd := &GitDaemon{
		repoPath:                            absoluteGitPath,
		watcher:                             w,
		debounceDur:                         debounce,
		gitFilesActiveRefreshDur:            gitFilesActiveRefreshDur,
		gitRemoteSyncStatusActiveRefreshDur: gitRemoteSyncStatusActiveRefreshDur,
		watcherTimer:                        time.NewTimer(debounce), // milliseconds
		gitFilesActiveTimer:                 time.NewTimer(gitFilesActiveRefreshDur),
		gitRemoteSyncStatusActiveTimer:      time.NewTimer(gitRemoteSyncStatusActiveRefreshDur),
		stopChannel:                         make(chan struct{}),
		errorLog:                            make([]error, 0),
		updateChannel:                       updateChannel,
		gitOperations:                       gitOperations,
	}
	gd.isGitFilesPassiveActiveRunning.Store(false)
	gd.isGitRemoteSyncStatusActiveRunning.Store(false)
	gd.isGitBranchPassiveRunning.Store(false)
	gd.isGitStashPassiveRunning.Store(false)
	gd.watcherTimer.Stop()
	gd.gitFilesActiveTimer.Stop()
	gd.gitRemoteSyncStatusActiveTimer.Stop()
	gd.watchPath()

	GITDAEMON = gd
}

func (gd *GitDaemon) watchPath() {
	err := gd.watcher.Add(gd.repoPath)
	if err != nil {
		gd.errorLog = append(gd.errorLog, err)
	}
	err = filepath.WalkDir(gd.repoPath, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			gd.watcher.Add(path)
		}
		return nil
	})
	if err != nil {
		gd.errorLog = append(gd.errorLog, err)
	}

}

func (gd *GitDaemon) Start() {
	go func() {
		// Initial call to get info of git
		if gd.updateChannel != nil {
			gd.gitLatestInfoFetch()
		}
		gd.gitFilesActiveTimer.Reset(gd.gitFilesActiveRefreshDur)
		gd.gitRemoteSyncStatusActiveTimer.Reset(gd.gitRemoteSyncStatusActiveRefreshDur)
		// loop to stay active
		for {
			select {
			case event := <-gd.watcher.Events:
				if gd.isRelevantEvent(event) {
					gd.resetDebounce()
				}
			case err := <-gd.watcher.Errors:
				fmt.Println("Watcher error:", err)

			case <-gd.watcherTimer.C:
				gd.gitLatestInfoFetch()
			case <-gd.gitFilesActiveTimer.C:
				// reset first to avoid losing ticks, then run work in background
				gd.gitFilesActiveTimer.Reset(gd.gitFilesActiveRefreshDur)
				go func() {
					if gd.isGitFilesPassiveActiveRunning.CompareAndSwap(false, true) {
						// Mark as running
						defer gd.isGitFilesPassiveActiveRunning.Store(false)

						gd.gitOperations.GitFiles.GetGitFilesStatus()
						gd.updateChannel <- git.GIT_FILES_STATUS_UPDATE
					}
				}()
			case <-gd.gitRemoteSyncStatusActiveTimer.C:
				// reset immediately; git remote sync status operation
				gd.gitRemoteSyncStatusActiveTimer.Reset(gd.gitRemoteSyncStatusActiveRefreshDur)
				go func() {
					if gd.isGitRemoteSyncStatusActiveRunning.CompareAndSwap(false, true) {
						defer gd.isGitRemoteSyncStatusActiveRunning.Store(false)
						gd.gitOperations.GitRemote.GetLatestRemoteSyncStatusAndUpstream()
						gd.updateChannel <- git.GIT_REMOTE_SYNC_STATUS_AND_UPSTREAM_UPDATE
					}
				}()
			case <-gd.stopChannel:
				gd.watcher.Close()
				return
			}
		}
	}()
}

func (gd *GitDaemon) resetDebounce() {
	if !gd.watcherTimer.Stop() {
		select {
		case <-gd.watcherTimer.C:
		default:
		}
	}
	gd.watcherTimer.Reset(gd.debounceDur)
}

func (gd *GitDaemon) gitLatestInfoFetch() {
	go func() {
		if gd.isGitFilesPassiveActiveRunning.CompareAndSwap(false, true) {
			defer gd.isGitFilesPassiveActiveRunning.Store(false)
			gd.gitOperations.GitFiles.GetGitFilesStatus()
			gd.updateChannel <- git.GIT_FILES_STATUS_UPDATE
		}
	}()
	go func() {
		if gd.isGitBranchPassiveRunning.CompareAndSwap(false, true) {
			defer gd.isGitBranchPassiveRunning.Store(false)
			gd.gitOperations.GitBranch.GetLatestBranchesinfo()
			gd.updateChannel <- git.GIT_BRANCH_UPDATE
		}
	}()
	go func() {
		if gd.isGitRemoteSyncStatusActiveRunning.CompareAndSwap(false, true) {
			defer gd.isGitRemoteSyncStatusActiveRunning.Store(false)
			gd.gitOperations.GitRemote.GetLatestRemoteSyncStatusAndUpstream()
			gd.updateChannel <- git.GIT_REMOTE_SYNC_STATUS_AND_UPSTREAM_UPDATE
		}
	}()
	go func() {
		if gd.isGitStashPassiveRunning.CompareAndSwap(false, true) {
			defer gd.isGitStashPassiveRunning.Store(false)
			gd.gitOperations.GitStash.GetLatestStashInfo()
			gd.updateChannel <- git.GIT_STASH_UPDATE
		}
	}()
}

func (gd *GitDaemon) isRelevantEvent(event fsnotify.Event) bool {
	// Only watch .git subpaths
	if !strings.Contains(event.Name, filepath.Join(gd.repoPath)) {
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
					_ = gd.watcher.Add(path)
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
	close(gd.stopChannel)
}
