package commitfiles

import "github.com/gohyuhan/gitti/tui/types"

func InitGitCheckoutCommitFileConfirmPromptPopUpModel(m *types.GittiModel, fileDisplayName string, filePathName string, commitHash string) {
	popUpModel := &GitCheckoutCommitFileConfirmPromptPopUpModel{
		FileDisplayName: fileDisplayName,
		FilePathName:    filePathName,
		CommitHash:      commitHash,
	}
	m.PopUpModel = popUpModel
}
