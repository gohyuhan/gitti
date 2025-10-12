package types

type BranchesInfo struct {
	Name            string
	CurrentCheckout bool
}

type FilesInfo struct {
	Name             string
	IncludedForStage bool
}

type GitInfo struct {
	CurrentCheckedOutBranch string
	AllBranches             map[string]BranchesInfo
	AllChangedFiles         map[string]FilesInfo
	CurrentSelectedFile     string
}
