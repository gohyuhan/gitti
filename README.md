# Gitti

A fast, lightweight terminal UI for Git operations that keeps you in your flow.

> ⚠️ **Development Status**: Gitti is actively under development. Features, APIs, and behaviors may change significantly. Development is driven by personal needs first, with additional features planned afterward. Not all majority used Git operations are currently supported.

## Description

Gitti is a visual Git client built for developers who live in the terminal. It provides an intuitive TUI (Terminal User Interface) for common Git operations without the overhead of traditional GUI applications or the verbosity of CLI commands.

## Why Gitti?
Gitti is built for terminal-focused developers who need visual Git operations without breaking their flow. Traditional GUI clients like GitHub Desktop offer great interfaces but consume significant RAM and force context switching that disrupts your coding rhythm. Pure CLI commands are powerful but lack visual context for reviewing changes and managing branches. Born from personal need while working in Neovim, Gitti bridges this gap by bringing an intuitive, lightweight TUI directly into your terminal, no window management, no context switching, just seamless Git operations with visual clarity. Plus, it's universal with built-in support for English, Japanese, Simplified & Traditional Chinese.

## Features

- 🌳 **Branch Management** - View, switch, and manage branches with ease
- 📝 **Interactive Staging** - Visually select and stage files
- 🔍 **Diff Viewer** - Review changes with syntax-aware diff display
- 💬 **Commit Interface** - Write commits with a dedicated UI
- 🚀 **Push/Pull Operations** - Manage remote operations seamlessly
- 💿 **Changes Stash Operations** - Manage stash operations seamlessly
- 📦 **Basic Submodule Support** - Work with Git submodules in your repositories
- 🌍 **Multi-language Support** - English, Japanese, 简体中文, 繁體中文
- ⚡ **Real-time Updates** - File system monitoring for instant status updates
- ⌨️ **Keyboard-driven** - Efficient navigation without touching the mouse

## Installation

### Linux
```bash
curl --proto "=https" -sSfL https://github.com/gohyuhan/gitti/releases/latest/download/install.sh | bash
```

### macOS (curl or homebrew)
```bash
curl --proto "=https" -sSfL https://github.com/gohyuhan/gitti/releases/latest/download/install.sh | bash

# via homebrew
# Add the tap (once)
brew tap gohyuhan/gitti

# Install latest
brew update && brew install gitti
```

### Windows (PowerShell or scoop)
```powershell
powershell -c "irm https://github.com/gohyuhan/gitti/releases/latest/download/install.ps1 | iex"

# via scoop
# Add the bucket (once)
scoop bucket add gitti https://github.com/gohyuhan/scoop-gitti

# Install latest
scoop update; scoop install gitti
```

### Go Install
If you have Go installed, you can install Gitti directly:
```bash
go install github.com/gohyuhan/gitti@latest
```

## Uninstall & Cleanup

### macOS (Homebrew)
```bash
# 1. Uninstall + remove ALL versions
brew uninstall --force gitti

# 2. Remove the tap
brew untap gohyuhan/gitti

# 3. Delete the binary directly (in case it's not a symlink or brew missed it)
rm -f /opt/homebrew/bin/gitti
rm -f /usr/local/bin/gitti

# 4. Delete the entire Cellar folder for gitti (old kegs)
rm -rf /opt/homebrew/Cellar/gitti
rm -rf /usr/local/Cellar/gitti

# 5. Delete any leftover symlinks
rm -rf /opt/homebrew/opt/gitti
rm -rf /usr/local/opt/gitti

# 6. Delete all cached downloads for gitti
rm -rf ~/Library/Caches/Homebrew/gitti*
rm -rf ~/Library/Caches/Homebrew/downloads/*gitti*
```

### Windows (Scoop)
```powershell
# 1. Uninstall the app (all versions)
scoop uninstall gitti 2>$null

# 2. Remove the bucket
scoop bucket rm gitti 2>$null

# 3. Delete the app folder completely (including shims + persist)
rm -r -force "$env:USERPROFILE\scoop\apps\gitti" 2>$null

# 4. Delete the bucket clone
rm -r -force "$env:USERPROFILE\scoop\buckets\gitti" 2>$null

# 5. Delete all cached installers for gitti
scoop cache rm "gitti*" 2>$null
```

### Manual Installation (curl / powershell)

#### macOS / Linux
```bash
# Remove binary (if installed via curl)
sudo rm -f /usr/local/bin/gitti
```

#### Windows
```powershell
# Remove binary and directory
Remove-Item -Path "$env:LOCALAPPDATA\gitti" -Recurse -Force
```

### Configuration Cleanup

To completely remove Gitti's configuration files:

#### macOS
```bash
rm -rf "$HOME/Library/Application Support/gitti"
```

#### Linux
```bash
rm -rf "$HOME/.config/gitti"
```

#### Windows
```powershell
Remove-Item -Path "$env:APPDATA\gitti" -Recurse -Force
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

### [v0.1.1]
- refactor: refactor bubble, bubble tea, lip gloss to be compatible for its upcoming v2 release
- Fix: Remove discard selection option
- Fix: Fix unstage all api to be compatible for unborn repo
- Fix: fix a slight ui issue where it will softwrap and break the layout
- refactor: update git operation api for create/switch branch to use native git command instead of using own custom logic flow to get the same result but more efficiently and faster

### [v0.1.0]
- Initial release in development
- Core TUI implementation
- Branch management and switching
- Interactive file staging
- Diff viewer
- Commit, pull and push operations
- Changes stash operation
- Multi-language support (en, ja, zh-hans, zh-hant)
- Real-time file system monitoring
- Configuration management
- Basic submodule support

## License

MIT License - see [LICENSE](LICENSE) file for details

## Acknowledgments

Built with:
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components

---

**Made with ❤️ for terminal enthusiasts who refuse to break their flow**
