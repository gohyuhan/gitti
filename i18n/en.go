package i18n

var eN = LanguageMapping{
	FlagLangCode:                       "set the language code, e.g. 'EN', 'JP', 'ZH-CN', 'ZH-TW'...",
	FlagInitDefaultBranch:              "set the default branch for git init.\nThis only affect the git init done through gitti.\nFor git to default to use the configure branch name please also add a '--global' flag",
	FlagGlobal:                         "also apply to system git if supported",
	FailToGetCWD:                       "Failed to get current working directory",
	TuiRunFail:                         "Alas, there's been an error",
	LanguageNotSupportedPanic:          "[%s] is not supported, Please choose between %v",
	LanguageSet:                        "Language set to %s",
	GittiDefaultBranchSet:              "Gitti default init branch set to '%s'",
	GittiDefaultAndGitDefaultBranchSet: "Both gitti and git default init branch set to '%s'",
	Branches:                           "Branches",
	ModifiedFiles:                      "Modified Files",
	FileTypeUnSupportedPreview:         "The current selected file type is not supported for preview",
	TerminalSizeWarning:                "Terminal too small — resize to continue.",
	CurrentTerminalHeight:              "Current height",
	MinimumTerminalHeight:              "Minimum required height",
	CurrentTerminalWidth:               "Current width",
	MinimumTerminalWidth:               "Minimum required height",
	KeyBindingNoneSelected: []string{
		"[b] branch component",
		"[f] files component",
		"[esc] quit",
	},
	KeyBindingLocalBranchComponentIsCheckOut: []string{
		"[s] stash all file(s)",
		"[u] unstage all file(s)",
		"[esc] unselect component",
	},
	KeyBindingLocalBranchComponentDefault: []string{
		"[enter] switch branch",
		"[s] stash all file(s)",
		"[u] unstage all file(s)",
		"[esc] unselect component",
	},
	KeyBindingLocalBranchComponentNone: []string{
		"[esc] unselect component",
	},
	KeyBindingModifiedFilesComponentIsStaged: []string{
		"[s] unstage this change",
		"[enter] view modified content",
		"[esc] unselect component",
	},
	KeyBindingModifiedFilesComponentDefault: []string{
		"[s] stage this change",
		"[enter] view modified content",
		"[esc] unselect component",
	},
	KeyBindingModifiedFilesComponentNone: []string{
		"[esc] unselect component",
	},
	KeyBindingFileDiffComponent: []string{
		"[←/→] move left and right",
		"[↑/↓] move up and down",
		"[esc] back to file component",
	},
}
