package i18n

var ZH_CN = LanguageMapping{
	Branches:                   "分支",
	ModifiedFiles:              "已修改的文件",
	FileTypeUnSupportedPreview: "当前选择的文件类型不支持预览",
	TerminalSizeWarning:        "终端窗口太小 — 请调整大小后继续。",
	CurrentTerminalHeight:      "当前高度",
	MinimumTerminalHeight:      "最小要求高度",
	CurrentTerminalWidth:       "当前宽度",
	MinimumTerminalWidth:       "最小要求宽度",
	KeyBindingNoneSelected: []string{
		"[b] 分支组件",
		"[f] 文件组件",
		"[esc] 退出",
	},
	KeyBindingLocalBranchComponentIsCheckOut: []string{
		"[s] 储藏所有文件（stash）",
		"[u] 取消所有文件的暂存（unstage）",
		"[esc] 取消选择组件",
	},
	KeyBindingLocalBranchComponentDefault: []string{
		"[enter] 切换分支",
		"[s] 储藏所有文件（stash）",
		"[u] 取消所有文件的暂存（unstage）",
		"[esc] 取消选择组件",
	},
	KeyBindingLocalBranchComponentNone: []string{
		"[esc] 取消选择组件",
	},
	KeyBindingModifiedFilesComponentIsStaged: []string{
		"[s] 取消暂存此更改",
		"[enter] 查看修改内容",
		"[esc] 取消选择组件",
	},
	KeyBindingModifiedFilesComponentDefault: []string{
		"[s] 暂存此更改",
		"[enter] 查看修改内容",
		"[esc] 取消选择组件",
	},
	KeyBindingModifiedFilesComponentNone: []string{
		"[esc] 取消选择组件",
	},
	KeyBindingFileDiffComponent: []string{
		"[←/→] 左右移动",
		"[↑/↓] 上下移动",
		"[esc] 返回文件组件",
	},
}
