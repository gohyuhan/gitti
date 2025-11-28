package git

import (
	"sync/atomic"

	"github.com/gohyuhan/gitti/i18n"
)

// this basically act as the universal lock so that any function that will exec git operation
// that involce any write that will trigger git internal lock will need to check and try to acquire this lock in gitti.
// If failed to acquire, the function should return directly to prevent any concurrent execution of git operation that involve write
type GitProcessLock struct {
	isGitLockedForProcessing      atomic.Bool
	otherGitProcessRunningWarning string
}

func InitGitProcessLock() *GitProcessLock {
	gPL := &GitProcessLock{
		otherGitProcessRunningWarning: i18n.LANGUAGEMAPPING.OtherGitOpsIsRunningWarning,
	}
	gPL.isGitLockedForProcessing.Store(false)

	return gPL
}

func (gpl *GitProcessLock) CanProceedWithGitOps() bool {
	if gpl.isGitLockedForProcessing.CompareAndSwap(false, true) {
		return true
	}

	return false
}

func (gpl *GitProcessLock) ReleaseGitOpsLock() {
	gpl.isGitLockedForProcessing.Store(false)
}

func (gpl *GitProcessLock) OtherProcessRunningWarning() string {
	return gpl.otherGitProcessRunningWarning
}
