package api

import (
	"fmt"
	"gitti/api/git"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

type GitDaemon struct {
	RepoPath                  string
	Watcher                   *fsnotify.Watcher
	DebounceDur               time.Duration
	GitFilesPassiveRefreshDur time.Duration
	GitFetchPassiveRefreshDur time.Duration
	Paused                    bool
	DebounceMU                sync.Mutex
	GitMU                     sync.Mutex
	WatcherTimer              *time.Timer
	GitFilesPassiveTimer      *time.Timer
	GitFetchPassiveTimer      *time.Timer
	StopChannel               chan struct{}
	ErrorLog                  []error
	UpdateChannel             chan string // to communicate back to main thread for an update event
}

var GITDAEMON *GitDaemon

func InitGitDaemon(repoPath string, debounce time.Duration, updateChannel chan string) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}

	gitFilesPassiveRefreshDur := 2 * time.Second
	gitFetchPassiveRefreshDur := 60 * time.Second
	gd := &GitDaemon{
		RepoPath:                  filepath.Join(repoPath, ".git"),
		Watcher:                   w,
		DebounceDur:               debounce,
		GitFilesPassiveRefreshDur: gitFilesPassiveRefreshDur,
		GitFetchPassiveRefreshDur: gitFetchPassiveRefreshDur,
		WatcherTimer:              time.NewTimer(debounce), // milliseconds
		GitFilesPassiveTimer:      time.NewTimer(gitFilesPassiveRefreshDur),
		GitFetchPassiveTimer:      time.NewTimer(gitFetchPassiveRefreshDur),
		StopChannel:               make(chan struct{}),
		ErrorLog:                  make([]error, 0),
		UpdateChannel:             updateChannel,
	}
	gd.WatcherTimer.Stop()
	gd.GitFilesPassiveTimer.Stop()
	gd.GitFetchPassiveTimer.Stop()
	gd.WatchPath()

	GITDAEMON = gd
	return
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
		gd.GitMU.Lock()
		if !gd.Paused && gd.UpdateChannel != nil {
			git.GetUpdatedGitInfo(gd.UpdateChannel)
		}
		gd.GitMU.Unlock()
		gd.GitFilesPassiveTimer.Reset(gd.GitFilesPassiveRefreshDur)
		gd.GitFetchPassiveTimer.Reset(gd.GitFetchPassiveRefreshDur)
		// loop to stay active
		for {
			select {
			case event := <-gd.Watcher.Events:
				if gd.isRelevantEvent(event) {
					go func() {
						gd.DebounceMU.Lock()
						if !gd.Paused {
							gd.resetDebounce()
						}
						gd.DebounceMU.Unlock()
					}()
				}
			case err := <-gd.Watcher.Errors:
				fmt.Println("Watcher error:", err)

			case <-gd.WatcherTimer.C:
				go func() {
					gd.GitMU.Lock()
					if !gd.Paused && gd.UpdateChannel != nil {
						git.GetUpdatedGitInfo(gd.UpdateChannel)
					}
					gd.GitMU.Unlock()
				}()
			case <-gd.GitFilesPassiveTimer.C:
				go func() {
					gd.GitMU.Lock()
					if !gd.Paused {
						git.GITFILES.GetGitFilesStatus()
						gd.UpdateChannel <- git.GENERAL_GIT_UPDATE
					}
					gd.GitFilesPassiveTimer.Reset(gd.GitFilesPassiveRefreshDur)
					gd.GitMU.Unlock()
				}()
			case <-gd.GitFetchPassiveTimer.C:
				go func() {
					gd.GitFetchPassiveTimer.Reset(gd.GitFetchPassiveRefreshDur)
				}()
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

func (gd *GitDaemon) Pause() {
	gd.GitMU.Lock()
	gd.Paused = true
	gd.GitMU.Unlock()
}

func (gd *GitDaemon) Resume() {
	gd.GitMU.Lock()
	gd.Paused = false
	gd.GitMU.Unlock()
}

func (gd *GitDaemon) Stop() {
	close(gd.StopChannel)
}
