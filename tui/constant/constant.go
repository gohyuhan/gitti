package constant

// -----------------------------------------------------------------------------
//
//	Pop Up Type
//
// -----------------------------------------------------------------------------
const (
	NoPopUp                                    = "NoPopUp"
	GlobalKeyBindingPopUp                      = "GlobalKeyBindingPopUp"
	AmendCommitPopUp                           = "AmendCommitPopUp"                           // IsTyping will be true
	CommitPopUp                                = "CommitPopUp"                                // IsTyping will be true
	AddRemotePromptPopUp                       = "AddRemotePromptPopUp"                       // IsTyping will be true
	ChoosePushTypePopUp                        = "ChoosePushTypePopUp"                        // IsTyping will be false
	ChooseRemotePopUp                          = "ChooseRemotePopUp"                          // IsTyping will be false
	GitRemotePushPopUp                         = "GitRemotePushPopUp"                         // IsTyping will be false
	ChooseNewBranchTypePopUp                   = "ChooseNewBranchTypePopUp"                   // IsTyping will be false
	CreateNewBranchPopUp                       = "CreateNewBranchPopUp"                       // IsTyping will be true
	ChooseSwitchBranchTypePopUp                = "ChooseSwitchBranchTypePopUp"                // IsTyping will be false
	SwitchBranchOutputPopUp                    = "SwitchBranchOutputPopUp"                    // IsTyping will be false
	ChooseGitPullTypePopUp                     = "ChooseGitPullTypePopUp"                     // IsTyping will be false
	GitPullOutputPopUp                         = "GitPullOutputPopUp"                         // IsTyping will be false
	GitStashMessagePopUp                       = "GitStashMessagePopUp"                       // IsTyping will be true
	GitDiscardTypeOptionPopUp                  = "GitDiscardTypeOptionPopUp"                  // IsTyping will be false
	GitDiscardConfirmPromptPopUp               = "GitDiscardConfirmPromptPopUp"               // IsTyping will be false
	GitStashOperationOutputPopUp               = "GitStashOperationOutputPopUp"               // IsTyping will be false
	GitStashConfirmPromptPopUp                 = "GitStashConfirmPromptPopUp"                 // IsTyping will be false
	GitResolveConflictOptionPopUp              = "GitResolveConflictOptionPopUp"              // IsTyping will be false
	GitDeleteBranchConfirmPromptPopUp          = "GitDeleteBranchConfirmPromptPopUp"          // IsTyping will be false
	GitDeleteBranchOutputPopUp                 = "GitDeleteBranchOutputPopUp"                 // IsTyping will be false
	CreateBranchBasedOnRemotePopUp             = "CreateBranchBasedOnRemotePopUp"             // IsTyping will be true
	CreateBranchBasedOnRemoteOutputPopUp       = "CreateBranchBasedOnRemoteOutputPopUp"       // IsTyping will be false
	GitResetLatestCommitTypeOptionPopUp        = "GitResetLatestCommitTypeOptionPopUp"        // IsTyping will be false
	GitResetLatestCommitConfirmPromptPopUp     = "GitResetLatestCommitConfirmPromptPopUp"     // IsTyping will be false
	GitResetToSelectedCommitTypeOptionPopUp    = "GitResetToSelectedCommitTypeOptionPopUp"    // IsTyping will be false
	GitResetToSelectedCommitConfirmPromptPopUp = "GitResetToSelectedCommitConfirmPromptPopUp" // IsTyping will be false
	GitCherryPickOptionSelectionPopUp          = "GitCherryPickOptionSelectionPopUp"          // IsTyping will be false
	GitCherryPickPopUp                         = "GitCherryPickPopUp"                         // IsTyping will be false
	GitEditCherryPickPopUp                     = "GitEditCherryPickPopUp"                     // IsTyping will be false
	GitCherryPickApplyConfirmPopUp             = "GitCherryPickApplyConfirmPopUp"             // IsTyping will be false
	GitDiscardFileLineChangeConfirmPopUp       = "GitDiscardFileLineChangeConfirmPopUp"       // IsTyping will be false
)

const (
	SelectedLeftPanelComponentHeightRatio = 0.40
	MinUnSelectedComponentPanelHeight     = 3 // this will ensure it will at least show 1 item within the list (so at its smallest. it will show, title, 1 item of the list and list counter)
)

const (
	MinWidth  = 80
	MinHeight = 24

	Padding                             = 1
	MainPageKeyBindingLayoutPanelHeight = 1

	MinLogComponentHeight   = 4
	LogComponentHeightRatio = 0.15

	ListItemOrTitleWidthPad = 4

	MaxGlobalKeyBindingPopUpWidth                      = 150
	MaxCommitPopUpWidth                                = 150
	MaxAmendCommitPopUpWidth                           = 150
	MaxAddRemotePromptPopUpWidth                       = 150
	MaxGitRemotePushPopUpWidth                         = 150
	MaxChooseRemotePopUpWidth                          = 150
	MaxChoosePushTypePopUpWidth                        = 150
	MaxChooseNewBranchTypePopUpWidth                   = 150
	MaxCreateNewBranchPopUpWidth                       = 150
	MaxChooseSwitchBranchTypePopUpWidth                = 150
	MaxSwitchBranchOutputPopUpWidth                    = 150
	MaxChooseGitPullTypePopUpWidth                     = 150
	MaxGitPullOutputPopUpWidth                         = 150
	MaxGitStashMessagePopUpWidth                       = 150
	MaxGitDiscardTypeOptionPopUpWidth                  = 150
	MaxGitDiscardConfirmPromptPopupWidth               = 150
	MaxGitStashOperationOutputPopUpWidth               = 150
	MaxGitStashConfirmPromptPopUpWidth                 = 150
	MaxGitResolveConflictOptionPopUpWidth              = 150
	MaxGitDeleteBranchConfirmPromptPopUpWidth          = 150
	MaxGitDeleteBranchOutputPopUpWidth                 = 150
	MaxCreateBranchBasedOnRemotePopUpWidth             = 150
	MaxCreateBranchBasedOnRemoteOutputPopUpWidth       = 150
	MaxGitResetLatestCommitTypeOptionPopUpWidth        = 150
	MaxGitResetLatestCommitConfirmPromptPopUpWidth     = 150
	MaxGitResetToSelectedCommitTypeOptionPopUpWidth    = 150
	MaxGitResetToSelectedCommitConfirmPromptPopUpWidth = 150
	MaxGitCherryPickOptionSelectionPopUpWidth          = 150
	MaxGitCherryPickPopUpWidth                         = 150
	MaxGitEditCherryPickPopUpWidth                     = 150
	MaxGitCherryPickApplyConfirmPopUpWidth             = 150
	MaxGitDiscardFileLineChangeConfirmPopUpWidth       = 150

	PopUpGlobalKeyBindingViewPortHeight                = 18
	PopUpGitCommitOutputViewPortHeight                 = 10
	PopUpGitAmendCommitOutputViewPortHeight            = 10
	PopUpAddRemoteOutputViewPortHeight                 = 2
	PopUpGitRemotePushOutputViewportHeight             = 10
	PopUpChooseRemoteHeight                            = 10
	PopUpChoosePushTypeHeight                          = 6
	PopUpChooseNewBranchTypeHeight                     = 6
	PopUpChooseSwitchBranchTypeHeight                  = 6
	PopUpSwitchBranchOutputViewPortHeight              = 10
	PopUpChooseGitPullTypeHeight                       = 6
	PopUpGitPullOutputViewportHeight                   = 16
	PopUpGitDiscardTypeOptionHeight                    = 6
	PopUpGitStashOperationOutputViewPortHeight         = 10
	PopUpGitResolveConflictOptionPopUpHeight           = 6
	PopUpGitDeleteBranchOutputViewportHeight           = 4
	PopUpCreateBranchBasedOnRemoteOutputViewportHeight = 4
	PopUpGitResetLatestCommitTypeOptionHeight          = 6
	PopUpGitResetToSelectedCommitTypeOptionHeight      = 6
	PopUpGitCherryPickOptionSelectionHeight            = 6
	PopUpGitCherryPickPopUpHeight                      = 10
	PopUpGitEditCherryPickPopUpHeight                  = 10
	PopUpGitDiscardFileLineChangeViewportHeight        = 1
)

// variables for indicating which panel/components/container or whatever the hell you wanna call it that the user is currently landed or selected, so that they can do precious action related to the part of whatever the hell you wanna call it
const (
	GitStatusComponent     = "C0" // component index 0
	LocalBranchComponent   = "C1" // component index 1
	ModifiedFilesComponent = "C2" // component index 2
	CommitLogComponent     = "C3" // component index 3
	StashComponent         = "C4" // component index 4

	LogComponent = "L0" // this can be selected by keybinding but not by number

	// this is not a selectable component from key binding but act like an extension for each component to enter for more detail,
	// no component index, the current selected component index will be still set as its parent's
	DetailComponent    = "EC-DT"  // extended component -  detail component
	DetailComponentTwo = "EC-DT2" // extended component -  detail component two (currently only used for unstaged changes diff)
)

// will be used by the key binding navigation of going to previous or next component panel
var ComponentNavigationList = []string{
	GitStatusComponent,
	LocalBranchComponent,
	ModifiedFilesComponent,
	CommitLogComponent,
	StashComponent,
}

const DETAIL_COMPONENT_PANEL_UPDATED = "DETAIL_COMPONENT_PANEL_UPDATED"

const (
	AUTHOR_GITHUB   = "https://github.com/gohyuhan"
	AUTHOR_LINKEDIN = "https://my.linkedin.com/in/yu-han-goh-209480200"
)

const (
	HORIZONTAL = "HORIZONTAL"
	VERTICAL   = "VERTICAL"
)

const (
	STAGE         = "STAGE"
	UNSTAGE       = "UNSTAGE"
	NOSTAGESTATUS = "NOSTAGESTATUS"
)

// action
const (
	PUSHACTION                = "PUSHACTION"
	CREATEBRANCHBASEDONREMOTE = "CREATEBRANCHBASEDONREMOTE"
)

// cherry pick ops options
const (
	CHERRYPICK      = "CHERRYPICK"
	EDITCHERRYPICK  = "EDITCHERRYPICK"
	APPLYCHERRYPICK = "APPLYCHERRYPICK"
)
