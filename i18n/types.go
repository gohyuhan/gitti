package i18n

// this was use to structure for the global keybinding
const (
	TITLE = "TITLE"
	INFO  = "INFO"
	WARN  = "WARN"
)

type KeyBindingMappingFormat struct {
	KeyBindingLine  string
	TitleOrInfoLine string
	LineType        string
}

type FeatureInstructionMappingFormat struct {
	Feature          string
	InstructionLines []string
	LineType         string
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
	FlagVersion               string
	FlagLangCode              string
	FlagInitDefaultBranch     string
	FlagAutoUpdate            string
	FlagUpdate                string
	FlagGlobal                string
	FlagEditor                string
	FlagMaxCommitLogCount     string
	FlagAllowCommitGraphWrite string
	FlagMaxLogCount           string
	FlagShowXLog              string
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
	GitNotInstalledError           string
	GitNotInitPrompt               string
	GitInitRefuse                  string
	GitInitPromptInvalidInput      string
	GitCertainStateStillInProgress string
	MaxCommitLogCountSet           string
	MaxCommitLogCountSetError      string
	AllowCommitGraphWriteEnabled   string
	AllowCommitGraphWriteDisabled  string
	AllowCommitGraphWriteSetError  string
	MaxLogCountSet                 string
	MaxLogCountSetError            string
	ShowXLogSet                    string
	ShowXLogSetError               string
	// Gitti UI text
	Branches                    string
	ModifiedFiles               string
	CommitLog                   string
	Stash                       string
	Tag                         string
	FileTypeUnSupportedPreview  string
	TerminalSizeWarning         string
	CurrentTerminalHeight       string
	MinimumTerminalHeight       string
	CurrentTerminalWidth        string
	MinimumTerminalWidth        string
	Loading                     string
	StagedTitle                 string
	UnstagedTitle               string
	LineEditingModeTitle        string
	CherryPickTitle             string
	EditCherryPickTitle         string
	ApplyCherryPickTitle        string
	CherryPickOpsSelectionTitle string
	CherryPickApplyConfirmTitle string
	// for Key Bindings
	KeyBindingForGitStatusComponent                         []string
	KeyBindingLocalBranchComponentIsCheckOut                []string
	KeyBindingLocalBranchComponentDefault                   []string
	KeyBindingLocalBranchComponentNone                      []string
	KeyBindingTagComponentNone                              []string
	KeyBindingTagComponentDefault                           []string
	KeyBindingModifiedFilesComponentConflict                []string
	KeyBindingModifiedFilesComponentIsStaged                []string
	KeyBindingModifiedFilesComponentDefault                 []string
	KeyBindingModifiedFilesComponentNone                    []string
	KeyBindingCommitLogComponent                            []string
	KeyBindingLogComponent                                  []string
	KeyBindingKeyDetailComponent                            []string
	KeyBindingKeyDetailComponentLineEditingEligible         []string
	KeyBindingKeyDetailComponentLineEditing                 []string
	KeyBindingKeyStashComponent                             []string
	KeyBindingKeyStashComponentNone                         []string
	KeyBindingForCommitPopUp                                []string
	KeyBindingForAmendCommitPopUp                           []string
	KeyBindingForAddRemotePromptPopUp                       []string
	KeyBindingForGitRemotePushPopUp                         []string
	KeyBindingForChooseRemotePopUp                          []string
	KeyBindingForChoosePushTypePopUp                        []string
	KeyBindingForChooseNewBranchTypePopUp                   []string
	KeyBindingForCreateNewBranchPopUp                       []string
	KeyBindingForChooseSwitchBranchTypePopUp                []string
	KeyBindingForSwitchBranchOutputPopUp                    []string
	KeyBindingForChooseGitPullTypePopUp                     []string
	KeyBindingForGitPullOutputPopUp                         []string
	KeyBindingForGitStashMessagePopUp                       []string
	KeyBindingForGitDiscardTypeOptionPopUp                  []string
	KeyBindingForGitDiscardConfirmPromptPopUp               []string
	KeyBindingForGitStashOperationOutputPopUp               []string
	KeyBindingForGitStashConfirmPromptPopUp                 []string
	KeyBindingForGitDeleteBranchOutputPopUp                 []string
	KeyBindingForGitDeleteBranchConfirmPromptPopUp          []string
	KeyBindingForCreateBranchBasedOnRemotePopUp             []string
	KeyBindingForCreateBranchBasedOnRemoteOutputPopUp       []string
	KeyBindingForGitResetLatestCommitTypeOptionPopUp        []string
	KeyBindingForGitResetLatestCommitConfirmPromptPopUp     []string
	KeyBindingForGitResetToSelectedCommitTypeOptionPopUp    []string
	KeyBindingForGitResetToSelectedCommitConfirmPromptPopUp []string
	KeyBindingForGitCherryPickOptionSelectionPopUp          []string
	KeyBindingForGitCherryPickPopUp                         []string
	KeyBindingForGitEditCherryPickPopUp                     []string
	KeyBindingForGitCherryPickApplyConfirmPopUp             []string
	KeyBindingForGitDiscardFileLineChangeConfirmPopUp       []string
	KeyBindingForKeybindingAndFeatureInstructionsPopUp      []string
	KeyBindingForCreateTagPopUp                             []string
	KeyBindingForCreateTagConfirmationPopUp                 []string
	KeyBindingForChooseDeleteTagOptionPopUp                 []string
	KeyBindingForChooseRemoteForDeleteRemoteTagPopUp        []string
	KeyBindingForDeleteTagOutputPopUp                       []string
	// -----------------
	//  For Pop Up
	// -----------------
	// Global Key KeyBinding
	GlobalKeyBinding []KeyBindingMappingFormat
	// Local Branch Component KeyBinding
	LocalBranchComponentKeyBinding []KeyBindingMappingFormat
	// Tag Component KeyBinding
	TagComponentKeyBinding []KeyBindingMappingFormat
	// Modified Files Component KeyBinding
	ModifiedFilesComponentKeyBinding []KeyBindingMappingFormat
	// Commit Log Component KeyBinding
	CommitLogComponentKeyBinding []KeyBindingMappingFormat
	// Stash Component KeyBinding
	StashComponentKeyBinding []KeyBindingMappingFormat
	// Log Component KeyBinding
	LogComponentKeyBinding []KeyBindingMappingFormat
	// Detail Component KeyBinding
	DetailComponentKeyBinding []KeyBindingMappingFormat
	// Feature Instructions
	FeatureInstructions []FeatureInstructionMappingFormat
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
	EnterRemoteBranchPrompt  string
	ChooseNewBranchTypeTitle string
	NewBranchInvalidWarning  string
	// Create Branch Option
	CreateNewBranchTitle                     string
	CreateNewBranchDescription               string
	CreateNewBranchAndSwitchTitle            string
	CreateNewBranchAndSwitchDescription      string
	CreateNewBranchBasedOnRemoteTitle        string
	CreateNewBranchBasedOnRemoteDescription  string
	RemoteOriginTitle                        string
	EnterRemoteBranchTitle                   string
	CreatingNewBranchBasedOnRemoteTitle      string
	CreatingNewBranchBasedOnRemoteProcessing string
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
	GitResolveConflictOptionTitle             string
	GitResolveConflictReset                   string
	GitResolveConflictAcceptOursChanges       string
	GitResolveConflictAcceptTheirsChanges     string
	GitResolveConflictResetInfo               string
	GitResolveConflictAcceptOursChangesInfo   string
	GitResolveConflictAcceptTheirsChangesInfo string
	// for git delete branch
	GitDeleteBranchTitle         string
	GitDeleteBranchComfirmPrompt string
	DeletingBranch               string
	// for git reset latest commit
	GitResetLatestCommitTypeOptionTitle       string
	GitResetToSelectedCommitTypeOptionTitle   string
	GitResetSoft                              string
	GitResetHard                              string
	GitResetMixed                             string
	GitResetSoftInfo                          string
	GitResetHardInfo                          string
	GitResetMixedInfo                         string
	GitResetLatestCommitSoftConfirmation      string
	GitResetLatestCommitHardConfirmation      string
	GitResetLatestCommitMixedConfirmation     string
	GitResetToSelectedCommitSoftInfo          string
	GitResetToSelectedCommitHardInfo          string
	GitResetToSelectedCommitMixedInfo         string
	GitResetToSelectedCommitSoftConfirmation  string
	GitResetToSelectedCommitHardConfirmation  string
	GitResetToSelectedCommitMixedConfirmation string
	// for cherry pick
	CherryPickOpsTitle            string
	CherryPickOpsDescription      string
	EditCherryPickOpsTitle        string
	EditCherryPickOpsDescription  string
	ApplyCherryPickOpsTitle       string
	ApplyCherryPickOpsDescription string
	CherryPickedFromBranch        string
	// for discard file line change
	GitDiscardFileLineChangeConfirmTitle string
	// for tag
	CreateTagPopUpNameTitle                 string
	CreateTagPopUpNameInputPlaceHolder      string
	CreateTagPopUpMessageTitle              string
	CreateTagPopUpMessageInputPlaceHolder   string
	CreateTagConfirmation                   string
	ChooseDeleteTagOptionTitle              string
	DeleteTagPopUpDeleteLocalTagOption      string
	DeleteTagPopUpDeleteLocalTagOptionInfo  string
	DeleteTagPopUpDeleteRemoteTagOption     string
	DeleteTagPopUpDeleteRemoteTagOptionInfo string
	DeleteTagOutputPopUpTitle               string
	DeleteTagDeleting                       string
}
