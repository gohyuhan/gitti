package i18n

var JP = LanguageMapping{
	Branches:                   "ブランチ",
	ModifiedFiles:              "変更されたファイル",
	FileTypeUnSupportedPreview: "現在選択されているファイル形式はプレビューに対応していません",
	TerminalSizeWarning:        "端末サイズが小さすぎます — サイズを変更してください。",
	CurrentTerminalHeight:      "現在の高さ",
	MinimumTerminalHeight:      "必要な最小の高さ",
	CurrentTerminalWidth:       "現在の幅",
	MinimumTerminalWidth:       "必要な最小の幅",
	KeyBindingNoneSelected: []string{
		"[b] ブランチコンポーネント",
		"[f] ファイルコンポーネント",
		"[esc] 終了",
	},
	KeyBindingLocalBranchComponentIsCheckOut: []string{
		"[s] すべてのファイルをスタッシュ",
		"[u] すべてのファイルをアンステージ",
		"[esc] コンポーネント選択を解除",
	},
	KeyBindingLocalBranchComponentDefault: []string{
		"[enter] ブランチを切り替え",
		"[s] すべてのファイルをスタッシュ",
		"[u] すべてのファイルをアンステージ",
		"[esc] コンポーネント選択を解除",
	},
	KeyBindingLocalBranchComponentNone: []string{
		"[esc] コンポーネント選択を解除",
	},
	KeyBindingModifiedFilesComponentIsStaged: []string{
		"[s] この変更をアンステージ",
		"[enter] 変更内容を表示",
		"[esc] コンポーネント選択を解除",
	},
	KeyBindingModifiedFilesComponentDefault: []string{
		"[s] この変更をステージ",
		"[enter] 変更内容を表示",
		"[esc] コンポーネント選択を解除",
	},
	KeyBindingModifiedFilesComponentNone: []string{
		"[esc] コンポーネント選択を解除",
	},
	KeyBindingFileDiffComponent: []string{
		"[←/→] 左右に移動",
		"[↑/↓] 上下に移動",
		"[esc] ファイルコンポーネントに戻る",
	},
}
