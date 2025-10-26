package i18n

var ZH_TW = LanguageMapping{
	Branches:                   "分支",
	ModifiedFiles:              "已修改的檔案",
	FileTypeUnSupportedPreview: "目前選擇的檔案類型不支援預覽",
	TerminalSizeWarning:        "終端機太小 — 請調整大小以繼續。",
	CurrentTerminalHeight:      "目前高度",
	MinimumTerminalHeight:      "最小需求高度",
	CurrentTerminalWidth:       "目前寬度",
	MinimumTerminalWidth:       "最小需求寬度",
	KeyBindingNoneSelected: []string{
		"[b] 分支元件",
		"[f] 檔案元件",
		"[esc] 離開",
	},
	KeyBindingLocalBranchComponentIsCheckOut: []string{
		"[s] 儲存所有變更（stash）",
		"[u] 取消所有檔案的暫存（unstage）",
		"[esc] 取消選取元件",
	},
	KeyBindingLocalBranchComponentDefault: []string{
		"[enter] 切換分支",
		"[s] 儲存所有變更（stash）",
		"[u] 取消所有檔案的暫存（unstage）",
		"[esc] 取消選取元件",
	},
	KeyBindingLocalBranchComponentNone: []string{
		"[esc] 取消選取元件",
	},
	KeyBindingModifiedFilesComponentIsStaged: []string{
		"[s] 取消暫存此變更",
		"[enter] 查看修改內容",
		"[esc] 取消選取元件",
	},
	KeyBindingModifiedFilesComponentDefault: []string{
		"[s] 暫存此變更",
		"[enter] 查看修改內容",
		"[esc] 取消選取元件",
	},
	KeyBindingModifiedFilesComponentNone: []string{
		"[esc] 取消選取元件",
	},
	KeyBindingFileDiffComponent: []string{
		"[←/→] 左右移動",
		"[↑/↓] 上下移動",
		"[esc] 返回檔案元件",
	},
}
