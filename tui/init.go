package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/v2/list"
	"github.com/charmbracelet/bubbles/v2/textarea"
	"github.com/charmbracelet/bubbles/v2/textinput"
	"github.com/charmbracelet/bubbles/v2/viewport"
	"gitti/api/git"
	"gitti/i18n"
) // this was for various components part init or reinit function due to update or newly create

func initBranchList(m *GittiModel) {
	items := []list.Item{
		gitBranchItem(git.GITBRANCH.CurrentCheckOut),
	}

	for _, branch := range git.GITBRANCH.AllBranches {
		items = append(items, gitBranchItem(branch))
	}

	m.CurrentRepoBranchesInfoList = list.New(items, gitBranchItemDelegate{}, m.HomeTabLeftPanelWidth, m.HomeTabLocalBranchesPanelHeight)
	m.CurrentRepoBranchesInfoList.Title = fmt.Sprintf("[b]  %s:", i18n.LANGUAGEMAPPING.Branches)
	m.CurrentRepoBranchesInfoList.SetShowStatusBar(false)
	m.CurrentRepoBranchesInfoList.SetFilteringEnabled(false)
	m.CurrentRepoBranchesInfoList.SetShowHelp(false)
	m.CurrentRepoBranchesInfoList.Styles.Title = titleStyle
	m.CurrentRepoBranchesInfoList.Styles.PaginationStyle = paginationStyle

	if m.NavigationIndexPosition.LocalBranchComponent > len(m.CurrentRepoBranchesInfoList.Items())-1 {
		m.CurrentRepoBranchesInfoList.Select(len(m.CurrentRepoBranchesInfoList.Items()) - 1)
	} else {
		m.CurrentRepoBranchesInfoList.Select(m.NavigationIndexPosition.LocalBranchComponent)
	}

	return
}

func initModifiedFilesList(m *GittiModel) {
	latestModifiedFilesArray := git.GITFILES.FilesStatus
	items := make([]list.Item, 0, len(latestModifiedFilesArray))
	for _, modifiedFile := range latestModifiedFilesArray {
		items = append(items, gitModifiedFilesItem(modifiedFile))
	}

	// get the previous selected file and see if it was within the new list if yes get the latest position of the previous selected file
	previousSelectedFile := m.CurrentRepoModifiedFilesInfoList.SelectedItem()
	selectedFilesPosition := -1

	for index, item := range items {
		if item == previousSelectedFile {
			selectedFilesPosition = index
			break
		}
	}

	m.CurrentRepoModifiedFilesInfoList = list.New(items, gitModifiedFilesItemDelegate{}, m.HomeTabLeftPanelWidth, m.HomeTabChangedFilesPanelHeight)
	m.CurrentRepoModifiedFilesInfoList.Title = fmt.Sprintf("[f] 📄%s:", i18n.LANGUAGEMAPPING.ModifiedFiles)
	m.CurrentRepoModifiedFilesInfoList.SetShowStatusBar(false)
	m.CurrentRepoModifiedFilesInfoList.SetFilteringEnabled(false)
	m.CurrentRepoModifiedFilesInfoList.SetShowHelp(false)
	m.CurrentRepoModifiedFilesInfoList.SetShowPagination(false)
	m.CurrentRepoModifiedFilesInfoList.Styles.Title = titleStyle

	if len(items) < 1 {
		return
	}

	if selectedFilesPosition >= 0 {
		m.CurrentRepoModifiedFilesInfoList.Select(selectedFilesPosition)
		m.NavigationIndexPosition.ModifiedFilesComponent = selectedFilesPosition
	} else {
		if m.NavigationIndexPosition.ModifiedFilesComponent > len(m.CurrentRepoModifiedFilesInfoList.Items())-1 {
			m.CurrentRepoModifiedFilesInfoList.Select(len(m.CurrentRepoModifiedFilesInfoList.Items()) - 1)
		} else {
			m.CurrentRepoModifiedFilesInfoList.Select(m.NavigationIndexPosition.ModifiedFilesComponent)
		}
	}
	return
}

// reinit and render diff file viewport
func reinitAndRenderModifiedFileDiffViewPort(m *GittiModel) {
	m.CurrentSelectedFileDiffViewportOffset = 0
	m.CurrentSelectedFileDiffViewport.SetXOffset(0)
	m.CurrentSelectedFileDiffViewport.SetYOffset(0)
	renderModifiedFilesDiffViewPort(m)
}

func initGitCommitPopUpModel(m *GittiModel) {
	CommitMessageTextInput := textinput.New()
	CommitMessageTextInput.Placeholder = i18n.LANGUAGEMAPPING.CommitPopUpMessageInputPlaceHolder
	CommitMessageTextInput.Focus()
	CommitMessageTextInput.VirtualCursor = true

	CommitDescriptionTextAreaInput := textarea.New()
	CommitDescriptionTextAreaInput.ShowLineNumbers = false
	CommitDescriptionTextAreaInput.Placeholder = i18n.LANGUAGEMAPPING.CommitPopUpCommitDescriptionInputPlaceHolder
	CommitDescriptionTextAreaInput.SetHeight(5)
	CommitDescriptionTextAreaInput.Blur()

	vp := viewport.New()
	vp.SoftWrap = true
	vp.MouseWheelEnabled = true
	vp.MouseWheelDelta = 1
	vp.SetHeight(popUpGitCommitOutputViewPortHeight)
	vp.SetWidth(min(maxCommitPopUpWidth, int(float64(m.Width)*0.8)) - 4)

	m.PopUpModel = &GitCommitPopUpModel{
		MessageTextInput:         CommitMessageTextInput,
		DescriptionTextAreaInput: CommitDescriptionTextAreaInput,
		TotalInputCount:          2,
		CurrentActiveInputIndex:  1,
		GitCommitOutputViewport:  vp,
		IsProcessing:             false,
		HasError:                 false,
		ProcessSuccess:           false,
	}
}
