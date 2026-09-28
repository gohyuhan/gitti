package nontyping

import (
	"testing"

	"charm.land/bubbles/v2/list"
	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/i18n"
	tagComponent "github.com/gohyuhan/gitti/tui/component/tag"
	"github.com/gohyuhan/gitti/tui/constant"
	branchPopUp "github.com/gohyuhan/gitti/tui/popup/branch"
	tagPopUp "github.com/gohyuhan/gitti/tui/popup/tag"
	"github.com/gohyuhan/gitti/tui/types"
)

func init() {
	i18n.InitGittiLanguageMapping("EN")
}

func setupTagModel(tagName string) *types.GittiModel {
	tagItems := []list.Item{
		tagComponent.GitTagItem{
			TagName: tagName,
		},
	}
	tagList := list.New(tagItems, tagComponent.GitTagItemDelegate{}, 80, 20)

	m := &types.GittiModel{
		Width:                                                     100,
		Height:                                                    40,
		CurrentSelectedComponent:                                  constant.LocalBranchOrTagOrRemoteOrWorktreeComponentPanel,
		CurrentLocalBranchOrTagOrRemoteOrWorktreeComponentShowing: constant.SHOW_TAG,
		CurrentRepoTagInfoList:                                    tagList,
		TuiUpdateChannel:                                          make(chan interface{}, 10),
		PopUpType:                                                 constant.NoPopUp,
	}
	m.ShowPopUp.Store(false)
	m.IsTyping.Store(false)
	return m
}

func TestTagEnter_OpensChooseTagActionPopUp(t *testing.T) {
	m := setupTagModel("v1.0.0")

	resModel, _ := handleNonTypingEnterKeyBindingInteraction(m)

	if !resModel.ShowPopUp.Load() {
		t.Fatalf("expected ShowPopUp to be true")
	}
	if resModel.PopUpType != constant.ChooseTagActionPopUp {
		t.Fatalf("expected PopUpType to be %s, got %s", constant.ChooseTagActionPopUp, resModel.PopUpType)
	}
	popUp, ok := resModel.PopUpModel.(*tagPopUp.ChooseTagActionPopUpModel)
	if !ok {
		t.Fatalf("expected PopUpModel to be *tagPopUp.ChooseTagActionPopUpModel")
	}
	if popUp.TagName != "v1.0.0" {
		t.Errorf("expected TagName 'v1.0.0', got '%s'", popUp.TagName)
	}
	if len(popUp.TagActionList.Items()) != 2 {
		t.Fatalf("expected 2 options, got %d", len(popUp.TagActionList.Items()))
	}
}

func TestTagEnter_SelectCheckoutTagOption(t *testing.T) {
	m := setupTagModel("v2.0.0")
	// Open choose tag action popup
	m, _ = handleNonTypingEnterKeyBindingInteraction(m)

	popUp := m.PopUpModel.(*tagPopUp.ChooseTagActionPopUpModel)
	popUp.TagActionList.Select(0) // Checkout option

	// Press enter on the option
	resModel, _ := handleNonTypingEnterKeyBindingInteraction(m)

	if resModel.PopUpType != constant.CheckoutTagOutputPopUp {
		t.Fatalf("expected PopUpType to be %s, got %s", constant.CheckoutTagOutputPopUp, resModel.PopUpType)
	}
	checkoutPopUp, ok := resModel.PopUpModel.(*tagPopUp.CheckoutTagOutputPopUpModel)
	if !ok {
		t.Fatalf("expected PopUpModel to be *tagPopUp.CheckoutTagOutputPopUpModel")
	}
	if checkoutPopUp.TagName != "v2.0.0" {
		t.Errorf("expected TagName 'v2.0.0', got '%s'", checkoutPopUp.TagName)
	}
}

func TestTagEnter_SelectCreateBranchOption(t *testing.T) {
	m := setupTagModel("v3.0.0")
	// Open choose tag action popup
	m, _ = handleNonTypingEnterKeyBindingInteraction(m)

	popUp := m.PopUpModel.(*tagPopUp.ChooseTagActionPopUpModel)
	popUp.TagActionList.Select(1) // Create branch option

	// Press enter on the option
	resModel, _ := handleNonTypingEnterKeyBindingInteraction(m)

	if resModel.PopUpType != constant.CreateNewBranchPopUp {
		t.Fatalf("expected PopUpType to be %s, got %s", constant.CreateNewBranchPopUp, resModel.PopUpType)
	}
	if !resModel.IsTyping.Load() {
		t.Errorf("expected IsTyping to be true")
	}
	branchPopUpModel, ok := resModel.PopUpModel.(*branchPopUp.CreateNewBranchPopUpModel)
	if !ok {
		t.Fatalf("expected PopUpModel to be *branchPopUp.CreateNewBranchPopUpModel")
	}
	if branchPopUpModel.CreateType != git.NEWBRANCHBASEDONTAG {
		t.Errorf("expected CreateType %s, got %s", git.NEWBRANCHBASEDONTAG, branchPopUpModel.CreateType)
	}
	if branchPopUpModel.BasedOnTagName != "v3.0.0" {
		t.Errorf("expected BasedOnTagName 'v3.0.0', got '%s'", branchPopUpModel.BasedOnTagName)
	}
}

func TestTagNKey_DirectlyOpensCreateBranchPopUp(t *testing.T) {
	m := setupTagModel("v4.0.0")

	resModel, _ := handleNonTypingnKeyBindingInteraction(m)

	if !resModel.ShowPopUp.Load() {
		t.Fatalf("expected ShowPopUp to be true")
	}
	if resModel.PopUpType != constant.CreateNewBranchPopUp {
		t.Fatalf("expected PopUpType to be %s, got %s", constant.CreateNewBranchPopUp, resModel.PopUpType)
	}
	if !resModel.IsTyping.Load() {
		t.Errorf("expected IsTyping to be true")
	}
	branchPopUpModel, ok := resModel.PopUpModel.(*branchPopUp.CreateNewBranchPopUpModel)
	if !ok {
		t.Fatalf("expected PopUpModel to be *branchPopUp.CreateNewBranchPopUpModel")
	}
	if branchPopUpModel.CreateType != git.NEWBRANCHBASEDONTAG {
		t.Errorf("expected CreateType %s, got %s", git.NEWBRANCHBASEDONTAG, branchPopUpModel.CreateType)
	}
	if branchPopUpModel.BasedOnTagName != "v4.0.0" {
		t.Errorf("expected BasedOnTagName 'v4.0.0', got '%s'", branchPopUpModel.BasedOnTagName)
	}
}
