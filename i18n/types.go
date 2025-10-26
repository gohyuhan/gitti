package i18n

// -------------------------------------------------------
//     EXAMPLE of the Language Data Structure
//     * the sequence and structure will follow EN's
// -------------------------------------------------------
//
// var EN = LanguageMapping{
// 	Branches:                   "Branches",
// 	ModifiedFiles:              "Modified Files",
// 	FileTypeUnSupportedPreview: "The current selected file type is not supported for preview",
// 	TerminalSizeWarning:        "Terminal too small — resize to continue.",
// 	CurrentTerminalHeight:      "Current height",
// 	MinimumTerminalHeight:      "Minimum required height",
// 	CurrentTerminalWidth:       "Current width",
// 	MinimumTerminalWidth:       "Minimum required height",
// 	KeyBindingNoneSelected: []string{
// 		"[b] branch component",
// 		"[f] files component",
// 		"[esc] quit",
// 	},
// 	KeyBindingLocalBranchComponentIsCheckOut: []string{
// 		"[s] stash all file(s)",
// 		"[u] unstage all file(s)",
// 		"[esc] unselect component",
// 	},
// 	KeyBindingLocalBranchComponentDefault: []string{
// 		"[enter] switch branch",
// 		"[s] stash all file(s)",
// 		"[u] unstage all file(s)",
// 		"[esc] unselect component",
// 	},
// 	KeyBindingLocalBranchComponentNone: []string{
// 		"[esc] unselect component",
// 	},
// 	KeyBindingModifiedFilesComponentIsStaged: []string{
// 		"[s] unstage this change",
// 		"[enter] view modified content",
// 		"[esc] unselect component",
// 	},
// 	KeyBindingModifiedFilesComponentDefault: []string{
// 		"[s] stage this change",
// 		"[enter] view modified content",
// 		"[esc] unselect component",
// 	},
// 	KeyBindingModifiedFilesComponentNone: []string{
// 		"[esc] unselect component",
// 	},
// 	KeyBindingFileDiffComponent: []string{
// 		"[←/→] move left and right",
// 		"[↑/↓] move up and down",
// 		"[esc] back to file component",
// 	},
// }

type LanguageMapping struct {
	Branches                   string
	ModifiedFiles              string
	FileTypeUnSupportedPreview string
	TerminalSizeWarning        string
	CurrentTerminalHeight      string
	MinimumTerminalHeight      string
	CurrentTerminalWidth       string
	MinimumTerminalWidth       string
	// for Key Bindings
	KeyBindingNoneSelected                   []string
	KeyBindingLocalBranchComponentIsCheckOut []string
	KeyBindingLocalBranchComponentDefault    []string
	KeyBindingLocalBranchComponentNone       []string
	KeyBindingModifiedFilesComponentIsStaged []string
	KeyBindingModifiedFilesComponentDefault  []string
	KeyBindingModifiedFilesComponentNone     []string
	KeyBindingFileDiffComponent              []string
}
