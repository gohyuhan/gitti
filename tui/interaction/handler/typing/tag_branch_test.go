package typing

import (
	"testing"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/constant"
	branchPopUp "github.com/gohyuhan/gitti/tui/popup/branch"
	"github.com/gohyuhan/gitti/tui/types"
)

func init() {
	i18n.InitGittiLanguageMapping("EN")
}

func TestTypingEnter_CreateNewBranchBasedOnTag(t *testing.T) {
	input := textinput.New()
	input.SetValue("new-tag-branch")

	popUpModel := &branchPopUp.CreateNewBranchPopUpModel{
		NewBranchNameInput: input,
		CreateType:         git.NEWBRANCHBASEDONTAG,
		BasedOnTagName:     "v1.0.0",
	}

	m := &types.GittiModel{
		PopUpType:  constant.CreateNewBranchPopUp,
		PopUpModel: popUpModel,
	}
	m.ShowPopUp.Store(true)
	m.IsTyping.Store(true)

	resModel, _ := handleTypingEnterKeyBindingInteraction(m, tea.KeyPressMsg{})

	if resModel.ShowPopUp.Load() {
		t.Errorf("expected ShowPopUp to be false after submission")
	}
	if resModel.IsTyping.Load() {
		t.Errorf("expected IsTyping to be false after submission")
	}
	if resModel.PopUpType != constant.NoPopUp {
		t.Errorf("expected PopUpType to be NoPopUp, got %s", resModel.PopUpType)
	}
	if resModel.PopUpModel != nil {
		t.Errorf("expected PopUpModel to be nil, got %v", resModel.PopUpModel)
	}
}
