package nontyping

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/logging"
	"github.com/gohyuhan/gitti/tui/component/branch"
	"github.com/gohyuhan/gitti/tui/component/commitfiles"
	"github.com/gohyuhan/gitti/tui/component/commitlog"
	"github.com/gohyuhan/gitti/tui/component/files"
	"github.com/gohyuhan/gitti/tui/component/reflog"
	"github.com/gohyuhan/gitti/tui/component/remote"
	"github.com/gohyuhan/gitti/tui/component/stash"
	"github.com/gohyuhan/gitti/tui/component/tag"
	"github.com/gohyuhan/gitti/tui/component/worktree"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/popup/copypopup"
	"github.com/gohyuhan/gitti/tui/types"
)

func handleNonTypingyKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if m.ShowPopUp.Load() {
		return m, nil
	}
	if m.CopyInProgress {
		m.GittiLogger.RegisterNewLog(logging.COPY_VALUE_OPS, "", logging.WARN, i18n.LANGUAGEMAPPING.CopyInProgress, false)
		return m, nil
	}
	labels := i18n.LANGUAGEMAPPING
	ready := func(label, value string) copypopup.Option {
		return copypopup.Option{Label: label, Value: value, Ready: true}
	}
	var options []copypopup.Option
	var lookupKind, lookupArgument string
	switch m.CurrentSelectedComponent {
	case constant.LocalBranchOrTagOrRemoteOrWorktreeComponentPanel:
		switch m.CurrentLocalBranchOrTagOrRemoteOrWorktreeComponentShowing {
		case constant.SHOW_LOCAL_BRANCH:
			item, ok := m.CurrentRepoBranchesInfoList.SelectedItem().(branch.GitBranchItem)
			if ok {
				options = []copypopup.Option{ready(labels.CopyBranchName, item.BranchName), {Label: labels.CopyUpstreamName}}
				lookupKind, lookupArgument = copypopup.Upstream, "refs/heads/"+item.BranchName
			}
		case constant.SHOW_TAG:
			item, ok := m.CurrentRepoTagInfoList.SelectedItem().(tag.GitTagItem)
			if ok {
				options = []copypopup.Option{ready(labels.CopyTagName, item.TagName)}
			}
		case constant.SHOW_REMOTE:
			item, ok := m.CurrentRepoRemoteInfoList.SelectedItem().(remote.GitRemoteItem)
			if ok {
				options = []copypopup.Option{ready(labels.CopyRemoteName, item.Name), ready(labels.CopyRemoteURL, item.Url)}
			}
		case constant.SHOW_WORKTREE:
			item, ok := m.CurrentRepoWorktreeInfoList.SelectedItem().(worktree.GitWorktreeItem)
			if ok {
				options = []copypopup.Option{ready(labels.CopyWorktreePath, item.WorktreePath)}
			}
		}
	case constant.ModifiedFilesComponentPanel:
		item, ok := m.CurrentRepoModifiedFilesInfoList.SelectedItem().(files.GitModifiedFilesItem)
		if ok {
			path := item.NewFilePathname
			options = []copypopup.Option{
				ready(labels.CopyRelativePath, path),
				ready(labels.CopyAbsolutePath, filepath.Join(m.RepoPath, path)),
				ready(labels.CopyFileName, filepath.Base(path)),
			}
		}
	case constant.CommitFilesComponentPanel:
		item, ok := m.CurrentRepoCommitFilesList.SelectedItem().(commitfiles.GitCommitFileItem)
		if ok {
			path := item.FilePathname
			options = []copypopup.Option{
				ready(labels.CopyRelativePath, path),
				ready(labels.CopyAbsolutePath, filepath.Join(m.RepoPath, path)),
				ready(labels.CopyFileName, filepath.Base(path)),
			}
		}
	case constant.CommitLogOrRefLogComponentPanel:
		switch m.CurrentCommitLogOrRefLogComponentShowing {
		case constant.SHOW_COMMITLOG:
			item, ok := m.CurrentRepoCommitLogInfoList.SelectedItem().(commitlog.GitCommitLogItem)
			if ok {
				options = []copypopup.Option{
					ready(labels.CopyFullHash, item.Hash),
					{Label: labels.CopyShortHash},
					ready(labels.CopySubject, item.Message),
					ready(labels.CopyAuthor, item.Author),
				}
				lookupKind, lookupArgument = copypopup.ShortHash, item.Hash
			}
		case constant.SHOW_REFLOG:
			item, ok := m.CurrentRepoRefLogInfoList.SelectedItem().(reflog.GitRefLogItem)
			if ok {
				options = []copypopup.Option{ready(labels.CopyFullHash, item.Hash), {Label: labels.CopyShortHash}}
				lookupKind, lookupArgument = copypopup.ShortHash, item.Hash
			}
		}
	case constant.StashComponentPanel:
		item, ok := m.CurrentRepoStashInfoList.SelectedItem().(stash.GitStashItem)
		if ok {
			options = []copypopup.Option{ready(labels.CopyStashRef, item.Id)}
		}
	}
	if len(options) == 0 {
		return m, nil
	}
	popupID := copypopup.Init(m, options)
	if lookupKind == "" {
		return m, nil
	}
	return m, copyLookupCmd(popupID, lookupKind, m.RepoPath, lookupArgument)
}

func copyLookupCmd(popupID uint64, kind, repoPath, argument string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		var args []string
		if kind == copypopup.ShortHash {
			args = []string{"rev-parse", "--short", argument}
		} else {
			args = []string{"for-each-ref", "--format=%(refname)%00%(upstream:short)", argument}
		}
		cmd := exec.CommandContext(ctx, "git", args...)
		cmd.Dir = repoPath
		output, err := cmd.CombinedOutput()
		if err != nil {
			err = fmt.Errorf("%w: %s", err, strings.TrimSpace(string(output)))
		}
		value := strings.TrimSpace(string(output))
		if kind == copypopup.Upstream && err == nil {
			value = ""
			for _, line := range strings.Split(string(output), "\n") {
				parts := strings.SplitN(line, "\x00", 2)
				if len(parts) == 2 && parts[0] == argument {
					value = strings.TrimSpace(parts[1])
					break
				}
			}
		}
		return types.CopyValueResolvedMsg{PopupID: popupID, Kind: kind, Value: value, Err: err}
	}
}
