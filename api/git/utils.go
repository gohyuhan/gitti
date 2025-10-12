package git

import (
	"fmt"
	"gitti/types"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
)

// get the Username and UserEmail for git operation
func GetGitUserAndEmail(repo *git.Repository) (string, string, error) {
	var userName, userEmail string

	// Try local repo config (.git/config)
	cfg, err := repo.Config()
	if err == nil && cfg.User.Name != "" && cfg.User.Email != "" {
		userName = cfg.User.Name
		userEmail = cfg.User.Email
	} else {
		// Try global config (~/.gitconfig)
		globalCfg, err := config.LoadConfig(config.GlobalScope)
		if err == nil && globalCfg.User.Name != "" && globalCfg.User.Email != "" {
			userName = globalCfg.User.Name
			userEmail = globalCfg.User.Email
		}
	}

	// Fallback to environment variables
	if userName == "" {
		userName = os.Getenv("GIT_AUTHOR_NAME")
	}
	if userEmail == "" {
		userEmail = os.Getenv("GIT_AUTHOR_EMAIL")
	}

	return userName, userEmail, nil
}

// setDefaultBranch sets the default branch for a repository, as go-git set to "master" as default
func setDefaultBranch(repo *git.Repository, branchName string) error {
	cfg, err := repo.Config()
	if err != nil {
		return err
	}

	// Set the initial branch in config
	cfg.Init.DefaultBranch = branchName

	// Set HEAD to point to the new branch
	headRef := plumbing.NewSymbolicReference(plumbing.HEAD, plumbing.NewBranchReferenceName(branchName))
	if err := repo.Storer.SetReference(headRef); err != nil {
		return err
	}

	// Remove the original master branch if it exists and we're not using "master" as our branch name
	if branchName != "master" {
		masterRef := plumbing.NewBranchReferenceName("master")
		if err := repo.Storer.RemoveReference(masterRef); err != nil {
			// It's okay if the reference doesn't exist, just log it
			fmt.Printf("Note: %v branch didn't exist or couldn't be removed: %v\n", masterRef, err)
		} else {
			fmt.Printf("Removed original master branch\n")
		}
	}

	return repo.SetConfig(cfg)
}

func GetCurrentBranch(repo *git.Repository) string {
	headRef, err := repo.Head()
	if err == nil {
		return headRef.Name().Short()
	}

	// Fallback: resolve symbolic HEAD reference (for uncommitted branches)
	ref, symErr := repo.Storer.Reference(plumbing.HEAD)
	if symErr == nil && ref.Target() != "" {
		return ref.Target().Short()
	}
	return ""
}

func GetAllBranches(repo *git.Repository, currentCheckedOutBranch string) map[string]types.BranchesInfo {
	branchesInfoArray := map[string]types.BranchesInfo{}
	branches, err := repo.Branches()
	if err != nil {
		return branchesInfoArray
	}

	_ = branches.ForEach(func(ref *plumbing.Reference) error {
		isCheckout := false
		branchName := ref.Name().Short()
		if branchName == currentCheckedOutBranch {
			isCheckout = true
		}
		branchesInfoArray[branchName] = types.BranchesInfo{
			Name:            branchName,
			CurrentCheckout: isCheckout,
		}

		return nil
	})

	return branchesInfoArray
}

func GetAllChangedFilesInWorkTree(repo *git.Repository) (map[string]string, string) {
	changedFiles := map[string]string{}
	var firstFile string
	wt, err := repo.Worktree()
	if err != nil {
		panic("Unusual Rare Error of can't retrieve worktree, you might wanna reinit the repo or continue with git cli")
	}

	status, err := wt.Status()
	if err != nil {
		panic("Unusual Rare Error of can't retrieve status of worktree, you might wanna reinit the repo or continue with git cli")
	}

	if len(status) == 0 {
		return changedFiles, ""
	}

	for file := range status {
		fileName := fmt.Sprintf("%v", file)
		changedFiles[fileName] = fileName

		if firstFile == "" {
			firstFile = fileName
		}
	}

	return changedFiles, firstFile
}
