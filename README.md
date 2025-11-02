# Gitti

A fast, lightweight terminal UI for Git operations that keeps you in your flow.

> ⚠️ **Development Status**: Gitti is actively under development. Features, APIs, and behaviors may change significantly. Development is driven by personal needs first, with additional features planned afterward. Not all Git operations are currently supported.

## Description

Gitti is a visual Git client built for developers who live in the terminal. It provides an intuitive TUI (Terminal User Interface) for common Git operations without the overhead of traditional GUI applications or the verbosity of CLI commands.

## Why Gitti?
Gitti is built for terminal-focused developers who need visual Git operations without breaking their flow. Traditional GUI clients like GitHub Desktop offer great interfaces but consume significant RAM and force context switching that disrupts your coding rhythm. Pure CLI commands are powerful but lack visual context for reviewing changes and managing branches. Born from personal need while working in Neovim, Gitti bridges this gap by bringing an intuitive, lightweight TUI directly into your terminal—no window management, no context switching, just seamless Git operations with visual clarity. Plus, it's universal with built-in support for English, Japanese, Simplified & Traditional Chinese.

## Features

- 🌳 **Branch Management** - View, switch, and manage branches with ease
- 📝 **Interactive Staging** - Visually select and stage files
- 🔍 **Diff Viewer** - Review changes with syntax-aware diff display
- 💬 **Commit Interface** - Write commits with a dedicated UI
- 🚀 **Push/Pull Operations** - Manage remote operations seamlessly
- 🌍 **Multi-language Support** - English, Japanese, 简体中文, 繁體中文
- ⚡ **Real-time Updates** - File system monitoring for instant status updates
- ⌨️ **Keyboard-driven** - Efficient navigation without touching the mouse

## Installation

> 🚧 Installation instructions will be available after the first release.

### From Binary (Coming Soon)
```bash
# Instructions will be added here
```

### From Source (Coming Soon)
```bash
# Instructions will be added here
```

## Quick Start

Launch Gitti in any Git repository:
```bash
gitti
```

### Configuration

Set your preferred language:
```bash
gitti --language en    # English
gitti --language ja    # Japanese
gitti --language zh-hans  # Simplified Chinese
gitti --language zh-hant  # Traditional Chinese
```

Configure default branch for new repositories:
```bash
# For gitti only
gitti --init-dbranch main

# For global Git configuration that will be set to gitti and system git
gitti --init-dbranch main --global
```

## Changelog

### [v0.1.0]
- Initial release in development
- Core TUI implementation
- Branch management and switching
- Interactive file staging
- Diff viewer
- Commit and push operations
- Multi-language support (en, ja, zh-hans, zh-hant)
- Real-time file system monitoring
- Configuration management

## License

MIT License - see [LICENSE](LICENSE) file for details

## Acknowledgments

Built with:
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components

---

**Made with ❤️ for terminal enthusiasts who refuse to break their flow**
