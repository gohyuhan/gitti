package types

type GitInfo struct {
	CurrentCheckedOutBranch string
	AllBranches             map[string]BranchesInfo
	AllChangedFiles         map[string]string
	CurrentSelectedFile     string
}
