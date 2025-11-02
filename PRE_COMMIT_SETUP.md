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

### 2. Install Go tools
```bash
# Install goimports (for import formatting and cleanup)
go install golang.org/x/tools/cmd/goimports@latest

# Install golangci-lint (comprehensive linter)
brew install golangci-lint
# or
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
```

## Setup

Install the pre-commit hooks:
```bash
pre-commit install
```

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
  - Formats import statements
  - Removes unused imports
  - Groups imports (stdlib, external, local)
  - Sorts imports alphabetically
- **go vet**: Runs Go's built-in static analyzer to find suspicious code
- **go mod tidy**: Cleans up `go.mod` and `go.sum` files
- **golangci-lint**: Runs comprehensive linting with multiple linters:
  - Code style and formatting
  - Error checking
  - Security issues (gosec)
  - Code complexity (gocyclo)
  - Duplicate code detection
  - Performance issues
  - And many more...

### Security
- **gitleaks**: Scans for secrets and credentials in code

## Manual Usage

Run all hooks on all files:
```bash
pre-commit run --all-files
```

Run specific hook:
```bash
pre-commit run golangci-lint --all-files
pre-commit run goimports --all-files
```

Skip hooks for a commit (not recommended):
```bash
git commit --no-verify
```

## Configuration Files

- `.pre-commit-config.yaml`: Pre-commit hook configuration
- `.golangci.yml`: golangci-lint configuration with enabled linters and rules

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
