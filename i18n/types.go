package i18n

// -------------------------------------------------------
//     Language Data Structure
//     * the sequence and structure will follow EN's
// -------------------------------------------------------

type LanguageMapping struct {
	// flag expalnation
	FlagLangCode          string
	FlagInitDefaultBranch string
	FlagGlobal            string
	// Run Error
	FailToGetCWD string
	TuiRunFail   string
	// i18n
	LanguageNotSupportedPanic string
	LanguageSet               string
	// init default branch
	GittiDefaultBranchSet              string
	GittiDefaultAndGitDefaultBranchSet string
	// Gitti terminal text
	GitNotInstalledError      string
	GitNotInitPrompt          string
	GitInitRefuse             string
	GitInitPromptInvalidInput string
	// Gitti UI text
	Branches                   string
	ModifiedFiles              string
	FileTypeUnSupportedPreview string
	TerminalSizeWarning        string
	CurrentTerminalHeight      string
	MinimumTerminalHeight      string
	CurrentTerminalWidth       string
	MinimumTerminalWidth       string
	// for Key Bindings
	KeyBindingNoneSelected                   []string
	KeyBindingLocalBranchComponentIsCheckOut []string
	KeyBindingLocalBranchComponentDefault    []string
	KeyBindingLocalBranchComponentNone       []string
	KeyBindingModifiedFilesComponentIsStaged []string
	KeyBindingModifiedFilesComponentDefault  []string
	KeyBindingModifiedFilesComponentNone     []string
	KeyBindingFileDiffComponent              []string
	KeyBindingForCommitPopUp                 []string
	KeyBindingForAddRemotePromptPopUp        []string
	KeyBindingForGitRemotePushPopUp          []string
	// -----------------
	//  For Pop Up
	// -----------------
	// commit
	CommitPopUpMessageTitle                      string
	CommitPopUpMessageInputPlaceHolder           string
	CommitPopUpDescriptionTitle                  string
	CommitPopUpCommitDescriptionInputPlaceHolder string
	// prompt to add remote origin
	AddRemotePrompt        string
	RemoteNameTitle        string
	RemoteNamePlaceHolder  string
	RemoteUrlTitle         string
	RemoteUrlPlaceHolder   string
	RemoteAddSuccess       string
	InvalidRemoteUrlFormat string
	// git push
	GitRemotePushTitle string
}
