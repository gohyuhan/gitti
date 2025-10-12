package git

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"gitti/types"
)

func GetGitInfo(repoPath string) (bool, int, error, types.GitInfo) {
	gitInitialInfo := types.GitInfo{CurrentSelectedFile: ""}

	// 1. Open the git repo path from the provided path
	repo, err := git.PlainOpen(repoPath)
	if err == git.ErrRepositoryNotExists {
		return false, repoNotInitYet, fmt.Errorf("Failed to get current repo that has been tracked by git: %v", err), gitInitialInfo
	}

	// 2. Get the current reference branch, if it was a newly init repo with one branch without any commit, get the symbolic HEAD reference
	currentCheckedOutBranch := GetCurrentBranch(repo)
	if currentCheckedOutBranch == "" {
		panic("Unusual Error for a repo that was initialized without any referencing branch, you might want to remove the .git folder and reinit again.")
	} else {
		gitInitialInfo.CurrentCheckedOutBranch = currentCheckedOutBranch
	}

	// 3. Get all branches, no branch info is expected for symbolic HEAD reference
	allBranches := GetAllBranches(repo, currentCheckedOutBranch)
	gitInitialInfo.AllBranches = allBranches

	// 4. Get all the files that are modified/had chnages
	allChangedFiles, firstFileEntry := GetAllChangedFilesInWorkTree(repo)
	gitInitialInfo.AllChangedFiles = allChangedFiles
	if len(allChangedFiles) > 0 {
		gitInitialInfo.CurrentSelectedFile = firstFileEntry
	}

	return true, getGitInfoSuccessfully, nil, gitInitialInfo
}

func GitInit(repoPath string, branchName string) (bool, int, error) {
	repo, initErr := git.PlainInit(repoPath, false)
	if initErr != nil {
		return false, failToInitRepo, fmt.Errorf("Failed Initialize Repo: %v", initErr)
	}

	if branchName != "master" {
		if setDefaultBranchErr := setDefaultBranch(repo, branchName); setDefaultBranchErr != nil {
			return false, failToSetDefaultBranchDuringInit, fmt.Errorf("Fail to set %s as default branch, current default branch is [master]", branchName)
		}
	}

	return true, initializeGitRepoSuccessfullt, nil
}
