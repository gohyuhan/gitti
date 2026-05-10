package constant

// -----------------------------------------------------------------------------
//
//	Pop Up Type
//
// -----------------------------------------------------------------------------
const (
	NoPopUp                                       = "NoPopUp"
	KeybindingAndFeatureInstructionsPopUp         = "KeybindingAndFeatureInstructionsPopUp"
	AmendCommitPopUp                              = "AmendCommitPopUp"                              // IsTyping will be true
	CommitPopUp                                   = "CommitPopUp"                                   // IsTyping will be true
	AddRemotePromptPopUp                          = "AddRemotePromptPopUp"                          // IsTyping will be true
	ChoosePushTypePopUp                           = "ChoosePushTypePopUp"                           // IsTyping will be false
	ChooseRemotePopUp                             = "ChooseRemotePopUp"                             // IsTyping will be false
	GitRemotePushPopUp                            = "GitRemotePushPopUp"                            // IsTyping will be false
	ChooseNewBranchTypePopUp                      = "ChooseNewBranchTypePopUp"                      // IsTyping will be false
	CreateNewBranchPopUp                          = "CreateNewBranchPopUp"                          // IsTyping will be true
	ChooseSwitchBranchTypePopUp                   = "ChooseSwitchBranchTypePopUp"                   // IsTyping will be false
	SwitchBranchOutputPopUp                       = "SwitchBranchOutputPopUp"                       // IsTyping will be false
	ChooseGitPullTypePopUp                        = "ChooseGitPullTypePopUp"                        // IsTyping will be false
	GitPullOutputPopUp                            = "GitPullOutputPopUp"                            // IsTyping will be false
	GitStashMessagePopUp                          = "GitStashMessagePopUp"                          // IsTyping will be true
	GitDiscardTypeOptionPopUp                     = "GitDiscardTypeOptionPopUp"                     // IsTyping will be false
	GitDiscardConfirmPromptPopUp                  = "GitDiscardConfirmPromptPopUp"                  // IsTyping will be false
	GitStashOperationOutputPopUp                  = "GitStashOperationOutputPopUp"                  // IsTyping will be false
	GitStashConfirmPromptPopUp                    = "GitStashConfirmPromptPopUp"                    // IsTyping will be false
	GitResolveConflictOptionPopUp                 = "GitResolveConflictOptionPopUp"                 // IsTyping will be false
	GitDeleteBranchConfirmPromptPopUp             = "GitDeleteBranchConfirmPromptPopUp"             // IsTyping will be false
	GitDeleteBranchOutputPopUp                    = "GitDeleteBranchOutputPopUp"                    // IsTyping will be false
	CreateBranchBasedOnRemotePopUp                = "CreateBranchBasedOnRemotePopUp"                // IsTyping will be true
	CreateBranchBasedOnRemoteOutputPopUp          = "CreateBranchBasedOnRemoteOutputPopUp"          // IsTyping will be false
	GitResetLatestCommitTypeOptionPopUp           = "GitResetLatestCommitTypeOptionPopUp"           // IsTyping will be false
	GitResetLatestCommitConfirmPromptPopUp        = "GitResetLatestCommitConfirmPromptPopUp"        // IsTyping will be false
	GitResetToSelectedCommitTypeOptionPopUp       = "GitResetToSelectedCommitTypeOptionPopUp"       // IsTyping will be false
	GitResetToSelectedCommitConfirmPromptPopUp    = "GitResetToSelectedCommitConfirmPromptPopUp"    // IsTyping will be false
	GitCherryPickOptionSelectionPopUp             = "GitCherryPickOptionSelectionPopUp"             // IsTyping will be false
	GitCherryPickPopUp                            = "GitCherryPickPopUp"                            // IsTyping will be false
	GitEditCherryPickPopUp                        = "GitEditCherryPickPopUp"                        // IsTyping will be false
	GitCherryPickApplyConfirmPopUp                = "GitCherryPickApplyConfirmPopUp"                // IsTyping will be false
	GitDiscardFileLineChangeConfirmPopUp          = "GitDiscardFileLineChangeConfirmPopUp"          // IsTyping will be false
	CreateTagPopUp                                = "CreateTagPopUp"                                // IsTyping will be true
	CreateTagConfirmationPopUp                    = "CreateTagConfirmationPopUp"                    // IsTyping will be false
	ChooseDeleteTagOptionPopUp                    = "ChooseDeleteTagOptionPopUp"                    // IsTyping will be false
	ChooseRemoteForDeleteRemoteTagPopUp           = "ChooseRemoteForDeleteRemoteTagPopUp"           // IsTyping will be false
	DeleteTagOutputPopUp                          = "DeleteTagOutputPopUp"                          // IsTyping will be false
	ChoosePushTagOptionPopUp                      = "ChoosePushTagOptionPopUp"                      // IsTyping will be false
	PushTagOutputPopUp                            = "PushTagOutputPopUp"                            // IsTyping will be false
	ChooseFetchTagOptionPopUp                     = "ChooseFetchTagOptionPopUp"                     // IsTyping will be false
	FetchTagOutputPopUp                           = "FetchTagOutputPopUp"                           // IsTyping will be false
	RemoveRemoteConfirmationPopUp                 = "RemoveRemoteConfirmationPopUp"                 // IsTyping will be false
	RemoteAsTrackingUpstreamConfirmationPopUp     = "RemoteAsTrackingUpstreamConfirmationPopUp"     // IsTyping will be false
	EditRemotePromptPopUp                         = "EditRemotePromptPopUp"                         // IsTyping will be true
	GitRevertParentOptionSelectionPopUp           = "GitRevertParentOptionSelectionPopUp"           // IsTyping will be false
	GitRevertConfirmationPopUp                    = "GitRevertConfirmationPopUp"                    // IsTyping will be false
	GitCherryPickFromRefLogApplyConfirmationPopUp = "GitCherryPickFromRefLogApplyConfirmationPopUp" // IsTyping will be false
	GitRebaseBranchInputPopUp                     = "GitRebaseBranchInputPopUp"                     // IsTyping will be true
	GitRebaseOutputPopUp                          = "GitRebaseOutputPopUp"                          // IsTyping will be false
	ChooseRemoteBranchOptionPopUp                 = "ChooseRemoteBranchOptionPopUp"                 // IsTyping will be false
	ChooseBranchOptionForMergePopUp               = "ChooseBranchOptionForMergePopUp"               // IsTyping will be false
	BranchMergeOutputPopUp                        = "BranchMergeOutputPopUp"                        // IsTyping will be false
	BlamePopUp                                    = "BlamePopUp"                                    // IsTyping will be true
	InteractiveRebaseOptionPopUp                  = "InteractiveRebaseOptionPopUp"                  // IsTyping will be false
	InteractiveRebaseFixupSquashSelectionPopUp    = "InteractiveRebaseFixupSquashSelectionPopUp"    // IsTyping will be false
	InteractiveRebaseFixupSquashCommitPopUp       = "InteractiveRebaseFixupSquashCommitPopUp"       // IsTyping will be true
	InteractiveRebaseFixupSquashOutputPopUp       = "InteractiveRebaseFixupSquashOutputPopUp"       // Istyping will be false
)

const (
	SelectedLeftPanelComponentHeightRatio = 0.40
	MinUnSelectedComponentPanelHeight     = 3 // this will ensure it will at least show 1 item within the list (so at its smallest. it will show, title, 1 item of the list and list counter)
)

const (
	MinWidth  = 80
	MinHeight = 24

	TextAreaInputMinHeight = 3
	TextAreaInputMaxHeight = 7

	Padding                             = 1
	MainPageKeyBindingLayoutPanelHeight = 1

	MinLogComponentHeight   = 4
	LogComponentHeightRatio = 0.15

	ListItemOrTitleWidthPad = 4

	MaxKeybindingAndFeatureInstructionsPopUpWidth         = 150
	MaxCommitPopUpWidth                                   = 150
	MaxAmendCommitPopUpWidth                              = 150
	MaxAddRemotePromptPopUpWidth                          = 150
	MaxGitRemotePushPopUpWidth                            = 150
	MaxChooseRemotePopUpWidth                             = 150
	MaxChoosePushTypePopUpWidth                           = 150
	MaxChooseNewBranchTypePopUpWidth                      = 150
	MaxCreateNewBranchPopUpWidth                          = 150
	MaxChooseSwitchBranchTypePopUpWidth                   = 150
	MaxSwitchBranchOutputPopUpWidth                       = 150
	MaxChooseGitPullTypePopUpWidth                        = 150
	MaxGitPullOutputPopUpWidth                            = 150
	MaxGitStashMessagePopUpWidth                          = 150
	MaxGitDiscardTypeOptionPopUpWidth                     = 150
	MaxGitDiscardConfirmPromptPopupWidth                  = 150
	MaxGitStashOperationOutputPopUpWidth                  = 150
	MaxGitStashConfirmPromptPopUpWidth                    = 150
	MaxGitResolveConflictOptionPopUpWidth                 = 150
	MaxGitDeleteBranchConfirmPromptPopUpWidth             = 150
	MaxGitDeleteBranchOutputPopUpWidth                    = 150
	MaxCreateBranchBasedOnRemotePopUpWidth                = 150
	MaxCreateBranchBasedOnRemoteOutputPopUpWidth          = 150
	MaxGitResetLatestCommitTypeOptionPopUpWidth           = 150
	MaxGitResetLatestCommitConfirmPromptPopUpWidth        = 150
	MaxGitResetToSelectedCommitTypeOptionPopUpWidth       = 150
	MaxGitResetToSelectedCommitConfirmPromptPopUpWidth    = 150
	MaxGitCherryPickOptionSelectionPopUpWidth             = 150
	MaxGitCherryPickPopUpWidth                            = 150
	MaxGitEditCherryPickPopUpWidth                        = 150
	MaxGitCherryPickApplyConfirmPopUpWidth                = 150
	MaxGitDiscardFileLineChangeConfirmPopUpWidth          = 150
	MaxCreateTagPopUpWidth                                = 150
	MaxCreateTagConfirmationPopUpWidth                    = 150
	MaxChooseDeleteTagOptionPopUpWidth                    = 150
	MaxChooseRemoteForDeleteRemoteTagPopUpWidth           = 150
	MaxDeleteTagOutputPopUpWidth                          = 150
	MaxChoosePushTagOptionPopUpWidth                      = 150
	MaxPushTagOutputPopUpWidth                            = 150
	MaxChooseFetchTagOptionPopUpWidth                     = 150
	MaxFetchTagOutputPopUpWidth                           = 150
	MaxRemoveRemoteConfirmationPopUpWidth                 = 150
	MaxRemoteAsTrackingUpstreamConfirmationPopUpWidth     = 150
	MaxEditRemotePromptPopUpWidth                         = 150
	MaxGitRevertParentOptionSelectionPopUpWidth           = 150
	MaxGitRevertConfirmationPopUpWidth                    = 150
	MaxGitCherryPickFromRefLogApplyConfirmationPopUpWidth = 150
	MaxGitRebaseBranchInputPopUpWidth                     = 150
	MaxGitRebaseOutputPopUpWidth                          = 150
	MaxChooseRemoteBranchOptionPopUpWidth                 = 150
	MaxChooseBranchOptionForMergePopUpWidth               = 150
	MaxBranchMergeOutputPopUpWidth                        = 150
	// BlamePopUpWidth will not be set as it will always take up 90% of width
	MaxInteractiveRebaseOptionPopUpWidth = 150
	// InteractiveRebaseFixupSquashSelectionPopUpWidth will not be set as it will always take up 90% of the width, 65% for selection, 35% for rebase display
	MaxInteractiveRebaseFixupSquashCommitPopUpWidth = 150
	MaxInteractiveRebaseFixupSquashOutputPopUpWidth = 150

	PopUpGlobalKeyBindingViewPortHeight                       = 18
	PopUpGitCommitOutputViewPortHeight                        = 10
	PopUpGitAmendCommitOutputViewPortHeight                   = 10
	PopUpAddRemoteOutputViewPortHeight                        = 2
	PopUpGitRemotePushOutputViewportHeight                    = 10
	PopUpChooseRemoteHeight                                   = 10
	PopUpChoosePushTypeHeight                                 = 8
	PopUpChooseNewBranchTypeHeight                            = 10
	PopUpChooseSwitchBranchTypeHeight                         = 6
	PopUpSwitchBranchOutputViewPortHeight                     = 10
	PopUpChooseGitPullTypeHeight                              = 8
	PopUpGitPullOutputViewportHeight                          = 16
	PopUpGitDiscardTypeOptionHeight                           = 8
	PopUpGitStashOperationOutputViewPortHeight                = 10
	PopUpGitResolveConflictOptionPopUpHeight                  = 8
	PopUpGitDeleteBranchOutputViewportHeight                  = 4
	PopUpCreateBranchBasedOnRemoteOutputViewportHeight        = 4
	PopUpGitResetLatestCommitTypeOptionHeight                 = 8
	PopUpGitResetToSelectedCommitTypeOptionHeight             = 8
	PopUpGitCherryPickOptionSelectionHeight                   = 8
	PopUpGitCherryPickPopUpHeight                             = 10
	PopUpGitEditCherryPickPopUpHeight                         = 10
	PopUpGitDiscardFileLineChangeViewportHeight               = 1
	PopUpChooseDeleteTagOptionHeight                          = 6
	PopUpChooseRemoteForDeleteRemoteTagHeight                 = 12
	PopUpDeleteTagOutputViewportHeight                        = 10
	PopUpChoosePushTagOptionHeight                            = 10
	PopUpPushTagOutputViewportHeight                          = 10
	PopUpChooseFetchTagOptionHeight                           = 10
	PopUpFetchTagOutputViewportHeight                         = 10
	PopUpGitRevertParentOptionSelectionHeight                 = 10
	PopUpGitRebaseOutputViewportHeight                        = 10
	PopUpChooseRemoteBranchOptionHeight                       = 10
	PopUpChooseBranchOptionForMergeBranchOptionHeight         = 5
	PopUpChooseBranchOptionForMergeSelectedBranchOptionHeight = 5
	PopUpBranchMergeOutputViewportHeight                      = 16
	// BlamePopUpWidth will not be set as it will always take up 90% of height
	PopUpInteractiveRebaseOptionHeight = 8
	// InteractiveRebaseFixupSquashSelectionPopUpWidth will not be set as it will always take up 90% of the width, 65% for selection, 35% for rebase display
	PopUpInteractiveRebaseFixupSquashOutputviewportHeight = 8
)

// variables for indicating which component panel or whatever the hell you wanna call it that the user is currently landed or selected, so that they can do precious action related to the part of whatever the hell you wanna call it
const (
	GitStatusComponentPanel                = "C0" // component panel index 0
	LocalBranchOrTagOrRemoteComponentPanel = "C1" // component panel index 1 (local branch component, tag component and remote component share the same panel)
	ModifiedFilesComponentPanel            = "C2" // component panel index 2
	CommitLogOrRefLogComponentPanel        = "C3" // component panel index 3 (commit log component and ref log component share the same panel)
	StashComponentPanel                    = "C4" // component panel index 4

	LogComponentPanel = "L0" // this can be selected by keybinding but not by number

	// this is not a selectable component from key binding but act like an extension for each component to enter for more detail,
	// no component index, the current selected component index will be still set as its parent's
	DetailComponentPanel    = "EC-DT"  // extended component panel -  detail component
	DetailComponentPanelTwo = "EC-DT2" // extended component panel -  detail component two (currently only used for unstaged changes diff)
)

// will be used by the key binding navigation of going to previous or next component panel
var ComponentPanelNavigationList = []string{
	GitStatusComponentPanel,
	LocalBranchOrTagOrRemoteComponentPanel,
	ModifiedFilesComponentPanel,
	CommitLogOrRefLogComponentPanel,
	StashComponentPanel,
}

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

// action that require remote origin
const (
	PUSHACTION                = "PUSHACTION"
	CREATEBRANCHBASEDONREMOTE = "CREATEBRANCHBASEDONREMOTE"
	TAGPUSHACTION             = "TAGPUSHACTION"
	TAGFETCHACTION            = "TAGFETCHACTION"
	REBASEACTION              = "REBASEACTION"
)

// cherry pick ops options
const (
	CHERRYPICK      = "CHERRYPICK"
	EDITCHERRYPICK  = "EDITCHERRYPICK"
	APPLYCHERRYPICK = "APPLYCHERRYPICK"
)

// to indicate which component is showing in LocalBranchOrTagOrRemoteComponentPanel
const (
	SHOW_LOCAL_BRANCH = "SHOW_LOCAL_BRANCH"
	SHOW_TAG          = "SHOW_TAG"
	SHOW_REMOTE       = "SHOW_REMOTE"
)

// to indicate which component is showing in CommitLogOrRefLogComponentPanel
const (
	SHOW_COMMITLOG = "SHOW_COMMITLOG"
	SHOW_REFLOG    = "SHOW_REFLOG"
)

// GITTI TUI UPDATE EVENT
const (
	DETAIL_COMPONENT_PANEL_LAYOUT_UPDATED_EVENT         = "DETAIL_COMPONENT_PANEL_LAYOUT_UPDATED_EVENT"
	DETAIL_COMPONENT_PANEL_LAYOUT_STATE_UPDATED_EVENT   = "DETAIL_COMPONENT_PANEL_LAYOUT_STATE_UPDATED_EVENT"
	DETAIL_COMPONENT_PANEL_LAYOUT_STATE_REINIT_EVENT    = "DETAIL_COMPONENT_PANEL_LAYOUT_STATE_REINIT_EVENT"
	GIT_SWITCH_BRANCH_RESULT_EVENT                      = "GIT_SWITCH_BRANCH_RESULT_EVENT"
	GIT_DELETE_BRANCH_RESULT_EVENT                      = "GIT_DELETE_BRANCH_RESULT_EVENT"
	GIT_CREATE_NEW_BRANCH_BASED_ON_REMOTE_RESULT_EVENT  = "GIT_CREATE_NEW_BRANCH_BASED_ON_REMOTE_RESULT_EVENT"
	GIT_CREATE_NEW_BRANCH_BASED_ON_REMOTE_INVALID_EVENT = "GIT_CREATE_NEW_BRANCH_BASED_ON_REMOTE_INVALID_EVENT"
	GIT_MERGE_RESULT_EVENT                              = "GIT_MERGE_RESULT_EVENT"
	GIT_DELETE_TAG_RESULT_EVENT                         = "GIT_DELETE_TAG_RESULT_EVENT"
	GIT_PUSH_TAG_RESULT_EVENT                           = "GIT_PUSH_TAG_RESULT_EVENT"
	GIT_FETCH_TAG_RESULT_EVENT                          = "GIT_FETCH_TAG_RESULT_EVENT"
	GIT_STASH_OPERATION_RESULT_EVENT                    = "GIT_STASH_OPERATION_RESULT_EVENT"
	GIT_ADD_REMOTE_RESULT_EVENT                         = "GIT_ADD_REMOTE_RESULT_EVENT"
	GIT_REBASE_RESULT_EVENT                             = "GIT_REBASE_RESULT_EVENT"
	GIT_PUSH_RESULT_EVENT                               = "GIT_PUSH_RESULT_EVENT"
	GIT_COMMIT_RESULT_EVENT                             = "GIT_COMMIT_RESULT_EVENT"
	GIT_AMEND_COMMIT_RESULT_EVENT                       = "GIT_AMEND_COMMIT_RESULT_EVENT"
	GIT_PULL_RESULT_EVENT                               = "GIT_PULL_RESULT_EVENT"
	INTERACTIVE_REBASE_FIXUP_SQUASH_RESULT_EVENT        = "INTERACTIVE_REBASE_FIXUP_SQUASH_RESULT_EVENT"
	INTERACTIVE_REBASE_FETCH_COMMITS_INFO_EVENT         = "INTERACTIVE_REBASE_FETCH_COMMITS_INFO_EVENT"
)
