package constant

// -----------------------------------------------------------------------------
//
//	Pop Up Type
//
// -----------------------------------------------------------------------------
const (
	NoPopUp                     = "NoPopUp"
	CommitPopUp                 = "CommitPopUp"                 // IsTyping will be true
	AddRemotePromptPopUp        = "AddRemotePromptPopUp"        // IsTyping will be true
	ChoosePushTypePopUp         = "ChoosePushTypePopUp"         // IsTyping will be false
	ChooseRemotePopUp           = "ChooseRemotePopUp"           // IsTyping will be false
	GitRemotePushPopUp          = "GitRemotePushPopUp"          // IsTyping will be false
	ChooseNewBranchTypePopUp    = "ChooseNewBranchTypePopUp"    // IsTyping will be false
	CreateNewBranchPopUp        = "CreateNewBranchPopUp"        // IsTyping will be true
	ChooseSwitchBranchTypePopUp = "ChooseSwitchBranchTypePopUp" // IsTyping will be false
	SwitchBranchOutputPopUp     = "SwitchBranchOutputPopUp"     // IsTyping will be false
	ChooseGitPullTypePopUp      = "ChooseGitPullTypePopUp"      // IsTyping will be false
	GitPullOutputPopUp          = "GitPullOutputPopUp"          // IsTyping will be false
)

const AUTOCLOSEINTERVAL = 500

const (
	MinWidth  = 80
	MinHeight = 24

	Padding                             = 1
	MainPageKeyBindingLayoutPanelHeight = 1

	ListItemOrTitleWidthPad = 5

	MaxLeftPanelWidth                   = 80
	MaxCommitPopUpWidth                 = 150
	MaxAddRemotePromptPopUpWidth        = 150
	MaxGitRemotePushPopUpWidth          = 150
	MaxChooseRemotePopUpWidth           = 150
	MaxChoosePushTypePopUpWidth         = 150
	MaxChooseNewBranchTypePopUpWidth    = 150
	MaxCreateNewBranchPopUpWidth        = 150
	MaxChooseSwitchBranchTypePopUpWidth = 150
	MaxSwitchBranchOutputPopUpWidth     = 150
	MaxChooseGitPullTypePopUpWidth      = 150
	MaxGitPullOutputPopUpWidth          = 150

	PopUpGitCommitOutputViewPortHeight     = 10
	PopUpAddRemoteOutputViewPortHeight     = 2
	PopUpGitRemotePushOutputViewportHeight = 10
	PopUpChooseRemoteHeight                = 10
	PopUpChoosePushTypeHeight              = 6
	PopUpChooseNewBranchTypeHeight         = 6
	PopUpChooseSwitchBranchTypeHeight      = 6
	PopUpSwitchBranchOutputViewPortHeight  = 10
	PopUpChooseGitPullTypeHeight           = 6
	PopUpGitPullOutputViewportHeight       = 10
)

// variables for indicating which panel/components/container or whatever the hell you wanna call it that the user is currently landed or selected, so that they can do precious action related to the part of whatever the hell you wanna call it
const (
	LocalBranchComponent   = "C1"
	ModifiedFilesComponent = "C2"
	FileDiffComponent      = "C3"
)
