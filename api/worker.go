package api

import (
	"fmt"
	"gitti/api/git"
	"gitti/types"
	"time"
)

const (
	defaultGitDaemonWorkerTick = 5
)

// goroutine worker thread to act as a daemon that can be start and stop to get git latest info of the repo to update the UI in real time
type GittiDaemonWorker struct {
	repoPath string
	updateCh chan types.GitInfo
	stopCh   chan struct{}
	running  bool
}

func NewGitWorkerDaemon(repoPath string) *GittiDaemonWorker {
	return &GittiDaemonWorker{
		repoPath: repoPath,
		updateCh: make(chan types.GitInfo, 3),
		stopCh:   make(chan struct{}),
	}
}

// Start launches the background goroutine
func (w *GittiDaemonWorker) Start() {
	if w.running {
		return
	}
	w.running = true

	go func() {
		defer func() { w.running = false }()

		ticker := time.NewTicker(defaultGitDaemonWorkerTick * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-w.stopCh:
				w.running = false
				return
			case <-ticker.C:
				go w.getGitInfo()
			}
		}
	}()
}

// Stop signals the goroutine to exit
func (w *GittiDaemonWorker) Stop() {
	if !w.running {
		return
	}
	close(w.stopCh)                // stop current worker
	w.stopCh = make(chan struct{}) // reset stop channel for next start
}

func (w *GittiDaemonWorker) ListenToUpdateChannel() chan types.GitInfo {
	return w.updateCh
}

func (w *GittiDaemonWorker) getGitInfo() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in getGitInfo:", r)
		}
	}()

	_, _, _, gitInfo := git.GetGitInfo(w.repoPath) // we assume it will always success and return something at this stage first
	w.updateCh <- gitInfo
}
