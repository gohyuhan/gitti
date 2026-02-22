package git

import (
	"sync/atomic"

	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/logging"
)

// this basically act as the universal lock so that any function that will exec git operation
// that involce any write that will trigger git internal lock will need to check and try to acquire this lock in gitti.
// If failed to acquire, the function should return directly to prevent any concurrent execution of git operation that involve write
type GitProcessLock struct {
	isGitLockedForProcessing      atomic.Bool
	otherGitProcessRunningWarning string
	logging                       *logging.GittiLogging
}

// ----------------------------------
//
//	Initialize the universal git process lock
//
// ----------------------------------
func InitGitProcessLock(logging *logging.GittiLogging) *GitProcessLock {
	gPL := &GitProcessLock{
		otherGitProcessRunningWarning: i18n.LANGUAGEMAPPING.OtherGitOpsIsRunningWarning,
		logging:                       logging,
	}
	gPL.isGitLockedForProcessing.Store(false)

	return gPL
}

// ----------------------------------
//
//	Attempt to acquire the git operations lock, returns true if successfully acquired
//
// ----------------------------------
func (gpl *GitProcessLock) CanProceedWithGitOps() bool {
	if gpl.isGitLockedForProcessing.CompareAndSwap(false, true) {
		return true
	} else {
		gpl.logging.RegisterNewLog(logging.ACQUIRE_GIT_LOCK, "", logging.WARN, "There is other git process that acquired a lock is still running", false)
		return false
	}
}

// ----------------------------------
//
//	Release the git operations lock
//
// ----------------------------------
func (gpl *GitProcessLock) ReleaseGitOpsLock() {
	gpl.isGitLockedForProcessing.Store(false)
}

// ----------------------------------
//
//	Return the warning message for when another git operation is already running
//
// ----------------------------------
func (gpl *GitProcessLock) OtherProcessRunningWarning() string {
	return gpl.otherGitProcessRunningWarning
}
