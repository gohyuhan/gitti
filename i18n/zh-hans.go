package i18n

var zH_HANS = LanguageMapping{
	FlagLangCode:                       "设置语言代码，例如: 'EN', 'JA', 'ZH-HANS', 'ZH-HANT'...",
	FlagInitDefaultBranch:              "设置 git init 的默认分支.\n此设置仅影响通过 gitti 执行的 git init.\n若要让 git 自身默认使用此分支名称，请同时添加 '--global' 参数.",
	FlagGlobal:                         "同时应用到系统 git（如果支持）",
	FailToGetCWD:                       "获取当前工作目录失败",
	TuiRunFail:                         "发生错误",
	LanguageNotSupportedPanic:          "不支持[%s]，请选择 %v 之一",
	LanguageSet:                        "语言已设置为 %s",
	GittiDefaultBranchSet:              "Gitti 默认初始化分支已设置为“%s”",
	GittiDefaultAndGitDefaultBranchSet: "Gitti 和 Git 的默认初始化分支已设置为“%s”",
	Branches:                           "分支",
	ModifiedFiles:                      "已修改的文件",
	FileTypeUnSupportedPreview:         "当前选择的文件类型不支持预览",
	TerminalSizeWarning:                "终端窗口太小 — 请调整大小后继续。",
	CurrentTerminalHeight:              "当前高度",
	MinimumTerminalHeight:              "最小要求高度",
	CurrentTerminalWidth:               "当前宽度",
	MinimumTerminalWidth:               "最小要求宽度",
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
