# Pre-commit Setup Guide

This project uses pre-commit hooks to maintain code quality and consistency.

## Prerequisites

### 1. Install pre-commit
```bash
# macOS
brew install pre-commit

# or using pip
pip install pre-commit
```

### 2. Ensure Go is installed
```bash
# Verify Go installation
go version

# The pre-commit hooks will automatically download goimports when needed
# No manual installation of Go tools is required
```

## Setup

Install the pre-commit hooks:
```bash
pre-commit install
```

On first commit, pre-commit will:
1. Download and cache the hook repositories
2. Download `goimports` tool (if not already cached)
3. Run all configured hooks

This initial setup may take a minute, but subsequent commits will be much faster.

## What the hooks do

### General Hooks
- **trailing-whitespace**: Removes trailing whitespace
- **end-of-file-fixer**: Ensures files end with a newline
- **check-yaml**: Validates YAML files
- **check-added-large-files**: Prevents committing large files
- **check-merge-conflict**: Detects merge conflict markers
- **check-case-conflict**: Detects case conflicts in filenames

### Go-specific Hooks
- **go fmt**: Formats Go code according to standard conventions
- **goimports**:
  - Automatically downloads and runs `goimports` tool
  - Formats import statements
  - Removes unused imports
  - Groups imports (stdlib, external, local packages)
  - Sorts imports alphabetically
  - Uses `-local gitti` flag to properly group local imports
- **go vet**: Runs Go's built-in static analyzer to find suspicious code
- **go mod tidy**: Cleans up `go.mod` and `go.sum` files (runs only when these files change)

### Security
- **gitleaks**: Scans for secrets and credentials in code

## Manual Usage

Run all hooks on all files:
```bash
pre-commit run --all-files
```

Run specific hook:
```bash
pre-commit run goimports --all-files
pre-commit run go-vet --all-files
pre-commit run gitleaks --all-files
```

Skip hooks for a commit (not recommended):
```bash
git commit --no-verify
```

## Configuration Files

- `.pre-commit-config.yaml`: Pre-commit hook configuration with all enabled hooks and their settings

## Current Hook Versions

- **pre-commit-hooks**: v4.6.0
- **gitleaks**: v8.19.1
- **goimports**: Latest version (auto-downloaded on first run)

## Troubleshooting

If hooks fail:
1. Read the error message carefully
2. Fix the issues manually or let the hooks auto-fix them
3. Stage the changes: `git add .`
4. Try committing again

To update hooks to latest versions:
```bash
pre-commit autoupdate
```
