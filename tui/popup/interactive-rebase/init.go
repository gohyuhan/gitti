package interactiverebase

import (
	"strings"
	"unicode/utf8"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
	"github.com/gohyuhan/gitti/tui/utils"
)

// ------------------------------------
//
//	Initialize the interactive rebase operation type selection popup, populating
//	a list with three options (fixup/squash, reword, drop). Filtering and
//	pagination are hidden; an item-count help key is attached.
//
// ------------------------------------
func InitInteractiveRebaseOptionPopUpModel(m *types.GittiModel) {
	interactiveRebaseOption := []InteractiveRebaseOptionItem{
		{
			Name:                  i18n.LANGUAGEMAPPING.InteractiveRebaseFixupSquash,
			Info:                  i18n.LANGUAGEMAPPING.InteractiveRebaseFixupSquashDescription,
			InteractiveRebaseType: git.FIXUPSQUASH,
		},
		{
			Name:                  i18n.LANGUAGEMAPPING.InteractiveRebaseReword,
			Info:                  i18n.LANGUAGEMAPPING.InteractiveRebaseFeatureComingSoon,
			InteractiveRebaseType: git.REWORD,
		},
		{
			Name:                  i18n.LANGUAGEMAPPING.InteractiveRebaseDrop,
			Info:                  i18n.LANGUAGEMAPPING.InteractiveRebaseFeatureComingSoon,
			InteractiveRebaseType: git.DROP,
		},
	}

	items := make([]list.Item, 0, len(interactiveRebaseOption))
	for _, interactiveRebaseOption := range interactiveRebaseOption {
		items = append(items, InteractiveRebaseOptionItem(interactiveRebaseOption))
	}

	width := (min(constant.MaxInteractiveRebaseOptionPopUpWidth, int(float64(m.Width)*0.8)) - 4)
	iROL := list.New(items, InteractiveRebaseOptionDelegate{}, width, constant.PopUpInteractiveRebaseOptionHeight)
	iROL.SetShowPagination(false)
	iROL.SetShowStatusBar(false)
	iROL.SetFilteringEnabled(false)
	iROL.SetShowTitle(false)

	// Custom Help Model for Count Display
	iROL.SetShowHelp(true)
	iROL.KeyMap = list.KeyMap{}
	iROL.Styles.HelpStyle = style.NewStyle.MarginTop(0).MarginBottom(0).PaddingTop(0).PaddingBottom(0)
	iROL.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &iROL, constant.MaxInteractiveRebaseOptionPopUpWidth)

	popUpModel := &InteractiveRebaseOptionPopUpModel{
		InteractiveRebaseOptionList: iROL,
	}

	m.PopUpModel = popUpModel
}

// ------------------------------------
//
//	Initialize the split-pane fixup/squash commit selection popup. The left pane
//	holds a commit list (initially empty); the right pane is a preview viewport.
//	A goroutine is spawned to fetch commit infos asynchronously and send them
//	via the TUI update channel once available.
//
// ------------------------------------
func InitInteractiveRebaseFixupSquashSelectionPopUpModel(m *types.GittiModel) {
	items := make([]list.Item, 0)
	popUpWidth := int(float64(m.Width) * 0.9)
	innerWidth := popUpWidth - 2
	listWidth := int(float64(innerWidth) * 0.65)
	vpWidth := innerWidth - listWidth

	selectedCommitHashMap := make(map[string]git.CommitInfo)
	height := int(float64(m.Height)*0.8) - 2
	iRFSSL := list.New(items, InteractiveRebaseFixupSquashSelectionDelegate{&selectedCommitHashMap}, listWidth-2, height)
	iRFSSL.SetShowPagination(false)
	iRFSSL.SetShowStatusBar(false)
	iRFSSL.SetFilteringEnabled(false)
	iRFSSL.SetShowTitle(false)

	// Custom Help Model for Count Display
	iRFSSL.SetShowHelp(true)
	iRFSSL.KeyMap = list.KeyMap{}
	iRFSSL.Styles.HelpStyle = style.NewStyle.MarginTop(0).MarginBottom(0).PaddingTop(0).PaddingBottom(0)
	iRFSSL.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &iRFSSL, m.Width)

	vp := viewport.New()
	vp.SoftWrap = true
	vp.MouseWheelEnabled = true
	vp.MouseWheelDelta = 1
	vp.SetHeight(height)
	vp.SetWidth(vpWidth - 2)

	vp.SetContent(i18n.LANGUAGEMAPPING.InteractiveRebaseFixupMustHaveAtLeastTwoSelectedError)

	popUpModel := &InteractiveRebaseFixupSquashSelectionPopUpModel{
		CommitList:                          iRFSSL,
		CommitFixupSquashViewport:           vp,
		OriginalRetrievedCommitList:         []git.CommitInfo{},
		SelectedCommitHashMap:               selectedCommitHashMap,
		IsCommitListSelected:                true,
		IsCommitFixupSquashViewportSelected: false,
	}

	go func() {
		commitInfos := m.GitOperations.GitInteractiveRebase.GetCommitInfos()
		listItems := make([]list.Item, 0, len(commitInfos))
		for _, commitInfo := range commitInfos {
			listItems = append(listItems, InteractiveRebaseFixupSquashSelectionItem{
				Hash:        commitInfo.Hash,
				Message:     commitInfo.Message,
				Author:      commitInfo.Author,
				Description: commitInfo.Description,
				Parent:      commitInfo.Parent,
				CommitOrder: commitInfo.CommitOrder,
			})
		}

		data := types.InteractiveRebaseFetchCommitInfoListEventDataInterface{
			PopUpModel:  constant.InteractiveRebaseFixupSquashSelectionPopUp,
			CommitInfos: commitInfos,
			ListItems:   listItems,
		}

		m.TuiUpdateChannel <- types.GittiTuiUpdateMsg{
			Event: constant.INTERACTIVE_REBASE_FETCH_COMMITS_INFO_EVENT,
			Data:  data,
		}
	}()

	m.PopUpModel = popUpModel
}

// ------------------------------------
//
//	Initialize the fixup/squash commit message editing popup. Pre-fills the
//	message and description inputs by constructing a combined commit message from
//	the sorted selected commits, with the oldest commit's message as the subject.
//	The message input starts focused; the description textarea starts blurred.
//
// ------------------------------------
func InitInteractiveRebaseFixupSquashCommitPopUp(m *types.GittiModel, originalRetrievedCommitList []git.CommitInfo, sortedSelectedCommits []git.CommitInfo) {
	message, description := constructInteractiveRebaseFixupSquashCommitMessageAndDescription(sortedSelectedCommits)
	commitMsg := message
	commitDesc := description
	commitMsgPlaceholder := i18n.LANGUAGEMAPPING.InteractiveRebaseFixupSquashCommitPopUpMessageInputPlaceHolder
	commitDescPlaceholder := i18n.LANGUAGEMAPPING.InteractiveRebaseFixupSquashCommitPopUpCommitDescriptionInputPlaceHolder

	CommitMessageTextInput := textinput.New()
	CommitMessageTextInput.SetValue(commitMsg)
	CommitMessageTextInput.Placeholder = commitMsgPlaceholder
	CommitMessageTextInput.Focus()
	CommitMessageTextInput.SetVirtualCursor(true)

	CommitDescriptionTextAreaInput := textarea.New()
	CommitDescriptionTextAreaInput.SetValue(commitDesc)
	CommitDescriptionTextAreaInput.ShowLineNumbers = false
	CommitDescriptionTextAreaInput.Placeholder = commitDescPlaceholder
	CommitDescriptionTextAreaInput.DynamicHeight = true
	CommitDescriptionTextAreaInput.MinHeight = constant.TextAreaInputMinHeight
	CommitDescriptionTextAreaInput.MaxHeight = constant.TextAreaInputMaxHeight
	CommitDescriptionTextAreaInput.MaxContentHeight = 9999
	CommitDescriptionTextAreaInput.MoveToEnd()
	CommitDescriptionTextAreaInput.Blur()

	popUpModel := &InteractiveRebaseFixupSquashCommitPopUpModel{
		MessageTextInput:            CommitMessageTextInput,
		DescriptionTextAreaInput:    CommitDescriptionTextAreaInput,
		TotalInputCount:             2,
		CurrentActiveInputIndex:     1,
		SortedSelectedCommits:       sortedSelectedCommits,
		OriginalRetrievedCommitList: originalRetrievedCommitList,
	}
	m.PopUpModel = popUpModel
}

// ------------------------------------
//
//	Build the initial commit message and description for a fixup/squash rebase.
//	Uses the oldest selected commit's message as the subject line and appends
//	bullet-style summaries of each subsequent commit to the description body.
//
// ------------------------------------
func constructInteractiveRebaseFixupSquashCommitMessageAndDescription(sortedSelectedCommits []git.CommitInfo) (string, string) {
	var message string
	var description strings.Builder

	message = sortedSelectedCommits[0].Message
	if utf8.RuneCountInString(strings.TrimSpace(sortedSelectedCommits[0].Description)) > 0 {
		description.WriteString(sortedSelectedCommits[0].Description)
		description.WriteRune('\n')
		description.WriteRune('\n')
	}

	for _, commit := range sortedSelectedCommits[1:] {
		description.WriteString("* " + commit.Message)
		description.WriteRune('\n')
		description.WriteRune('\n')
		if utf8.RuneCountInString(strings.TrimSpace(commit.Description)) > 0 {
			description.WriteString(commit.Description)
			description.WriteRune('\n')
			description.WriteRune('\n')
		}
	}

	return message, description.String()
}

// ------------------------------------
//
//	Initialize the fixup/squash rebase output popup with a scrollable viewport
//	and a spinner. All atomic state flags (IsProcessing, IsCancelled, HasError,
//	ProcessSuccess) are reset to false.
//
// ------------------------------------
func InitInteractiveRebaseFixupSquashOutputPopUpModel(m *types.GittiModel) {
	vp := viewport.New()
	vp.SoftWrap = true
	vp.MouseWheelEnabled = true
	vp.MouseWheelDelta = 1
	vp.SetHeight(constant.PopUpGitRebaseOutputViewportHeight)
	vp.SetWidth(min(constant.MaxGitRebaseOutputPopUpWidth, int(float64(m.Width)*0.8)) - 4)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = style.SpinnerStyle

	popUpModel := &InteractiveRebaseFixupSquashOutputPopUpModel{
		FixupSquashOutputViewport: vp,
		Spinner:                   s,
	}

	popUpModel.IsProcessing.Store(false)
	popUpModel.IsCancelled.Store(false)
	popUpModel.HasError.Store(false)
	popUpModel.ProcessSuccess.Store(false)

	m.PopUpModel = popUpModel
}
