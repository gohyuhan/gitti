package i18n

// this was use to structure for the global keybinding
const (
	TITLE = "TITLE"
	INFO  = "INFO"
	WARN  = "WARN"
)

type GlobalKeyBindingMappingFormat struct {
	KeyBindingLine  string
	TitleOrInfoLine string
	LineType        string
}

// -------------------------------------------------------
//
//	Language Data Structure
//	* the sequence and structure will follow EN's
//
// -------------------------------------------------------
type LanguageMapping struct {
	// flag expalnation
	FlagLangCode                string
	FlagInitDefaultBranch       string
	FlagInitStashGitIgnoredFile string
	FlagGlobal                  string
	// Run Error
	FailToGetCWD                string
	TuiRunFail                  string
	OtherGitOpsIsRunningWarning string
	// i18n
	LanguageNotSupportedPanic string
	LanguageSet               string
	// init default branch
	GittiDefaultBranchSet              string
	GittiDefaultAndGitDefaultBranchSet string
	// update to stash ignored file setting
	GittiStashIgnoredFile        string
	GittiWillNotStashIgnoredFile string
	// Gitti terminal text
	GitNotInstalledError      string
	GitNotInitPrompt          string
	GitInitRefuse             string
	GitInitPromptInvalidInput string
	// Gitti UI text
	Branches                   string
	ModifiedFiles              string
	Stash                      string
	FileTypeUnSupportedPreview string
	TerminalSizeWarning        string
	CurrentTerminalHeight      string
	MinimumTerminalHeight      string
	CurrentTerminalWidth       string
	MinimumTerminalWidth       string
	// for Key Bindings
	KeyBindingForGittiStatusComponent        []string
	KeyBindingLocalBranchComponentIsCheckOut []string
	KeyBindingLocalBranchComponentDefault    []string
	KeyBindingLocalBranchComponentNone       []string
	KeyBindingModifiedFilesComponentIsStaged []string
	KeyBindingModifiedFilesComponentDefault  []string
	KeyBindingModifiedFilesComponentNone     []string
	KeyBindingKeyDetailComponent             []string
	KeyBindingKeyStashComponent              []string
	KeyBindingKeyStashComponentNone          []string
	KeyBindingForCommitPopUp                 []string
	KeyBindingForAmendCommitPopUp            []string
	KeyBindingForAddRemotePromptPopUp        []string
	KeyBindingForGitRemotePushPopUp          []string
	KeyBindingForChooseRemotePopUp           []string
	KeyBindingForChoosePushTypePopUp         []string
	KeyBindingForChooseNewBranchTypePopUp    []string
	KeyBindingForCreateNewBranchPopUp        []string
	KeyBindingForChooseSwitchBranchTypePopUp []string
	KeyBindingForSwitchBranchOutputPopUp     []string
	KeyBindingForChooseGitPullTypePopUp      []string
	KeyBindingForGitPullOutputPopUp          []string
	KeyBindingForGitStashMessagePopUp        []string
	KeyBindingForGlobalKeyBindingPopUp       []string
	// -----------------
	//  For Pop Up
	// -----------------
	// Global Key KeyBinding
	GlobalKeyBinding []GlobalKeyBindingMappingFormat
	// commit
	CommitPopUpMessageTitle                                  string
	CommitPopUpMessageInputPlaceHolder                       string
	CommitPopUpDescriptionTitle                              string
	CommitPopUpCommitDescriptionInputPlaceHolder             string
	CommitPopUpProcessing                                    string
	CommitPopUpMessageTitleAmendVersion                      string
	CommitPopUpMessageInputPlaceHolderAmendVersion           string
	CommitPopUpDescriptionTitleAmendVersion                  string
	CommitPopUpCommitDescriptionInputPlaceHolderAmendVersion string
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
	// switch branch
	ChooseSwitchBranchTypeTitle string
	// Switch Branch Option
	SwitchBranchTitle                  string
	SwitchBranchDescription            string
	SwitchBranchWithChangesTitle       string
	SwitchBranchWithChangesDescription string
	// for switch branch output
	SwitchBranchSwitchingToPopUpTitle            string
	SwitchBranchPopUpSwitchProcessing            string
	SwitchBranchPopUpSwitchWithChangesProcessing string
	// Git Pull Option
	ChoosePullOptionPrompt string
	GitPullOption          string
	GitPullRebaseOption    string
	GitPullMergeOption     string
	// for git pull output
	GitPullTitle      string
	GitPullProcessing string
	// for stash message prompt
	GitStashMessageTitle       string
	GitStashMessagePlaceholder string
}
