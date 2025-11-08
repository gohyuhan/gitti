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
	KeyBindingForChooseRemotePopUp           []string
	KeyBindingForChoosePushTypePopUp         []string
	KeyBindingForChooseNewBranchTypePopUp    []string
	KeyBindingForCreateNewBranchPopUp        []string
	// -----------------
	//  For Pop Up
	// -----------------
	// commit
	CommitPopUpMessageTitle                      string
	CommitPopUpMessageInputPlaceHolder           string
	CommitPopUpDescriptionTitle                  string
	CommitPopUpCommitDescriptionInputPlaceHolder string
	CommitPopUpProcessing                        string
	// prompt to add remote origin
	AddRemotePopUpPrompt                 string
	AddRemotePopUpRemoteNameTitle        string
	AddRemotePopUpRemoteNamePlaceHolder  string
	AddRemotePopUpRemoteUrlTitle         string
	AddRemotePopUpRemoteUrlPlaceHolder   string
	AddRemotePopUpRemoteAddSuccess       string
	AddRemotePopUpInvalidRemoteUrlFormat string
	// git push
	GitRemotePushPopUpTitle      string
	GitRemotePushPopUpProcessing string
	GitRemotePushOptionTitle     string
	// Choose Remote
	ChooseRemoteTitle string
	// Choose push option
	NormalPush         string
	ForcePushSafe      string
	ForcePushDangerous string
	// Create New Branch
	CreateNewBranchPrompt    string
	ChooseNewBranchTypeTitle string
	NewBranchInvalidWarning  string
	// Create Branch Option
	CreateNewBranchTitle                string
	CreateNewBranchDescription          string
	CreateNewBranchAndSwitchTitle       string
	CreateNewBranchAndSwitchDescription string
}
