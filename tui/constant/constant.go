package constant

// -----------------------------------------------------------------------------
//
//	Pop Up Type
//
// -----------------------------------------------------------------------------
const (
	NoPopUp                      = "NoPopUp"
	GlobalKeyBindingPopUp        = "GlobalKeyBindingPopUp"
	AmendCommitPopUp             = "AmendCommitPopUp"             // IsTyping will be true
	CommitPopUp                  = "CommitPopUp"                  // IsTyping will be true
	AddRemotePromptPopUp         = "AddRemotePromptPopUp"         // IsTyping will be true
	ChoosePushTypePopUp          = "ChoosePushTypePopUp"          // IsTyping will be false
	ChooseRemotePopUp            = "ChooseRemotePopUp"            // IsTyping will be false
	GitRemotePushPopUp           = "GitRemotePushPopUp"           // IsTyping will be false
	ChooseNewBranchTypePopUp     = "ChooseNewBranchTypePopUp"     // IsTyping will be false
	CreateNewBranchPopUp         = "CreateNewBranchPopUp"         // IsTyping will be true
	ChooseSwitchBranchTypePopUp  = "ChooseSwitchBranchTypePopUp"  // IsTyping will be false
	SwitchBranchOutputPopUp      = "SwitchBranchOutputPopUp"      // IsTyping will be false
	ChooseGitPullTypePopUp       = "ChooseGitPullTypePopUp"       // IsTyping will be false
	GitPullOutputPopUp           = "GitPullOutputPopUp"           // IsTyping will be false
	GitStashMessagePopUp         = "GitStashMessagePopUp"         // IsTyping will be true
	GitDiscardTypeOptionPopUp    = "GitDiscardTypeOptionPopUp"    // IsTyping will be false
	GitDiscardConfirmPromptPopup = "GitDiscardConfirmPromptPopup" // IsTyping will be false
)

const AUTOCLOSEINTERVAL = 500

const SelectedLeftPanelComponentHeightRatio = 0.4

const (
	MinWidth  = 80
	MinHeight = 24

	Padding                             = 1
	MainPageKeyBindingLayoutPanelHeight = 1

	ListItemOrTitleWidthPad = 4

	MaxLeftPanelWidth                    = 80
	MaxGlobalKeyBindingPopUpWidth        = 150
	MaxCommitPopUpWidth                  = 150
	MaxAmendCommitPopUpWidth             = 150
	MaxAddRemotePromptPopUpWidth         = 150
	MaxGitRemotePushPopUpWidth           = 150
	MaxChooseRemotePopUpWidth            = 150
	MaxChoosePushTypePopUpWidth          = 150
	MaxChooseNewBranchTypePopUpWidth     = 150
	MaxCreateNewBranchPopUpWidth         = 150
	MaxChooseSwitchBranchTypePopUpWidth  = 150
	MaxSwitchBranchOutputPopUpWidth      = 150
	MaxChooseGitPullTypePopUpWidth       = 150
	MaxGitPullOutputPopUpWidth           = 150
	MaxGitStashMessagePopUpWidth         = 150
	MaxGitDiscardTypeOptionPopUpWidth    = 150
	MaxGitDiscardConfirmPromptPopupWidth = 150

	PopUpGlobalKeyBindingViewPortHeight     = 30
	PopUpGitCommitOutputViewPortHeight      = 10
	PopUpGitAmendCommitOutputViewPortHeight = 10
	PopUpAddRemoteOutputViewPortHeight      = 2
	PopUpGitRemotePushOutputViewportHeight  = 10
	PopUpChooseRemoteHeight                 = 10
	PopUpChoosePushTypeHeight               = 6
	PopUpChooseNewBranchTypeHeight          = 6
	PopUpChooseSwitchBranchTypeHeight       = 6
	PopUpSwitchBranchOutputViewPortHeight   = 10
	PopUpChooseGitPullTypeHeight            = 6
	PopUpGitPullOutputViewportHeight        = 16
	PopUpGitDiscardTypeOptionHeight         = 6
)

// variables for indicating which panel/components/container or whatever the hell you wanna call it that the user is currently landed or selected, so that they can do precious action related to the part of whatever the hell you wanna call it
const (
	GitStatusComponent     = "C0" // component index 0
	LocalBranchComponent   = "C1" // component index 1
	ModifiedFilesComponent = "C2" // component index 2
	StashComponent         = "C3" // component index 3

	// this is not a selectable component from key binding but act like an extension for each component to enter for more detail,
	// no component index, the current selected component index will be still set as its parent's
	DetailComponent = "EC-DT" // extended component -  detail component
)

// will be used by the key binding navigation of going to previous or next component panel
var ComponentNavigationList = []string{
	GitStatusComponent,
	LocalBranchComponent,
	ModifiedFilesComponent,
	StashComponent,
}

const DETAIL_COMPONENT_PANEL_UPDATED = "DETAIL_COMPONENT_PANEL_UPDATED"
