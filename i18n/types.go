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
	AboutGitti []string
	// Updater related
	UpdaterDownloadPrompt               string
	UpdaterAlreadyLatest                string
	UpdaterFailToCheckForUpdate         string
	UpdaterUnSupportedOS                string
	UpdaterDownloadFail                 string
	UpdaterBinaryReplaceFail            string
	UpdaterDownloading                  string
	UpdaterDownloadUnexpectedStatusCode string
	UpdaterDownloadSuccess              string
	UpdaterRequiresSudo                 string
	UpdaterAutoUpdaterEnable            string
	UpdaterAutoUpdaterDisable           string
	UpdaterAutoUpdaterSetError          string
	// flag expalnation
	FlagVersion           string
	FlagLangCode          string
	FlagInitDefaultBranch string
	FlagAutoUpdate        string
	FlagUpdate            string
	FlagGlobal            string
	FlagEditor            string
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
	// set editor related
	EditorTitle       string
	EditorDescription string
	EditorInstruction string
	EditorSetError    string
	EditorSetSuccess  string
	// Gitti terminal text
	GitNotInstalledError      string
	GitNotInitPrompt          string
	GitInitRefuse             string
	GitInitPromptInvalidInput string
	// Gitti UI text
	Branches                   string
	ModifiedFiles              string
	CommitLog                  string
	Stash                      string
	FileTypeUnSupportedPreview string
	TerminalSizeWarning        string
	CurrentTerminalHeight      string
	MinimumTerminalHeight      string
	CurrentTerminalWidth       string
	MinimumTerminalWidth       string
	Loading                    string
	StagedTitle                string
	UnstagedTitle              string
	// for Key Bindings
	KeyBindingForGitStatusComponent                []string
	KeyBindingLocalBranchComponentIsCheckOut       []string
	KeyBindingLocalBranchComponentDefault          []string
	KeyBindingLocalBranchComponentNone             []string
	KeyBindingModifiedFilesComponentConflict       []string
	KeyBindingModifiedFilesComponentIsStaged       []string
	KeyBindingModifiedFilesComponentDefault        []string
	KeyBindingModifiedFilesComponentNone           []string
	KeyBindingCommitLogComponent                   []string
	KeyBindingKeyDetailComponent                   []string
	KeyBindingKeyStashComponent                    []string
	KeyBindingKeyStashComponentNone                []string
	KeyBindingForCommitPopUp                       []string
	KeyBindingForAmendCommitPopUp                  []string
	KeyBindingForAddRemotePromptPopUp              []string
	KeyBindingForGitRemotePushPopUp                []string
	KeyBindingForChooseRemotePopUp                 []string
	KeyBindingForChoosePushTypePopUp               []string
	KeyBindingForChooseNewBranchTypePopUp          []string
	KeyBindingForCreateNewBranchPopUp              []string
	KeyBindingForChooseSwitchBranchTypePopUp       []string
	KeyBindingForSwitchBranchOutputPopUp           []string
	KeyBindingForChooseGitPullTypePopUp            []string
	KeyBindingForGitPullOutputPopUp                []string
	KeyBindingForGitStashMessagePopUp              []string
	KeyBindingForGitDiscardTypeOptionPopUp         []string
	KeyBindingForGitDiscardConfirmPromptPopUp      []string
	KeyBindingForGitStashOperationOutputPopUp      []string
	KeyBindingForGitStashConfirmPromptPopUp        []string
	KeyBindingForGitDeleteBranchOutputPopUp        []string
	KeyBindingForGitDeleteBranchConfirmPromptPopUp []string
	KeyBindingForGlobalKeyBindingPopUp             []string
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
	// for git discard type option list
	GitDiscardTypeOptionTitle     string
	GitDiscardWhole               string
	GitDiscardUnstage             string
	GitDiscardAndRevertRename     string
	GitDiscardWholeInfo           string
	GitDiscardUnstageInfo         string
	GitDiscardAndRevertRenameInfo string
	// for discard confirmation prompt
	GitDiscardWholeConfirmation            string
	GitDiscardUnstageConfirmation          string
	GitDiscardUntrackedConfirmation        string
	GitDiscardNewlyAddedorCopyConfirmation string
	GitDiscardAndRevertRenameConfirmation  string
	// for stash operation title (used in output pop up)
	GitStashAllTitle   string
	GitStashFileTitle  string
	GitStashApplyTitle string
	GitStashDropTitle  string
	GitStashPopTitle   string
	// for stash operation processing (used in output pop up)
	GitStashAllProcessing   string
	GitStashFileProcessing  string
	GitStashApplyProcessing string
	GitStashDropProcessing  string
	GitStashPopProcessing   string
	// for stash operation confirm prompt
	GitStashAllConfirmation   string
	GitStashFileConfirmation  string
	GitApplyStashConfirmation string
	GitDropStashConfirmation  string
	GitPopStashConfirmation   string
	// for resolve conflict option list
	GitResolveConflictOptionTitle               string
	GitResolveConflictReset                     string
	GitResolveConflictAcceptLocalChanges        string
	GitResolveConflictAcceptIncomingChanges     string
	GitResolveConflictResetInfo                 string
	GitResolveConflictAcceptLocalChangesInfo    string
	GitResolveConflictAcceptIncomingChangesInfo string
	// for git delete branch
	GitDeleteBranchTitle         string
	GitDeleteBranchComfirmPrompt string
	DeletingBranch               string
}
