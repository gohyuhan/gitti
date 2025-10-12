package types

type GitInitialInfo struct {
	CurrentCheckedOutBranch string
	AllBranches             []BranchesInfo
	AllChangedFiles         []string
	CurrentSelectedFile     string
}
