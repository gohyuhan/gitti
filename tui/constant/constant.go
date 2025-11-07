package constant

// -----------------------------------------------------------------------------
//
//	Pop Up Type
//
// -----------------------------------------------------------------------------
const (
	NoPopUp              = "NoPopUp"
	CommitPopUp          = "CommitPopUp"          // IsTyping will be true
	AddRemotePromptPopUp = "AddRemotePromptPopUp" // IsTyping will be true
	ChoosePushTypePopUp  = "ChoosePushTypePopUp"  // IsTyping will be false
	ChooseRemotePopUp    = "ChooseRemotePopUp"    // IsTyping will be false
	GitRemotePushPopUp   = "GitRemotePushPopUp"   // IsTyping will be false
)

const AUTOCLOSEINTERVAL = 500

const (
	MinWidth  = 80
	MinHeight = 24

	Padding                             = 1
	MainPageKeyBindingLayoutPanelHeight = 1

	ListItemOrTitleWidthPad = 5

	MaxLeftPanelWidth            = 80
	MaxCommitPopUpWidth          = 100
	MaxAddRemotePromptPopUpWidth = 100
	MaxGitRemotePushPopUpWidth   = 100
	MaxChooseRemotePopUpWidth    = 100
	MaxChoosePushTypePopUpWidth  = 100

	PopUpGitCommitOutputViewPortHeight     = 10
	PopUpAddRemoteOutputViewPortHeight     = 2
	PopUpGitRemotePushOutputViewportHeight = 10
	PopUpChooseRemoteHeight                = 10
	PopUpChoosePushTypeHeight              = 10
)

// variables for indicating which panel/components/container or whatever the hell you wanna call it that the user is currently landed or selected, so that they can do precious action related to the part of whatever the hell you wanna call it
const (
	NoneSelected = "0"

	LocalBranchComponent   = "C1"
	ModifiedFilesComponent = "C2"
	FileDiffComponent      = "C3"
)
