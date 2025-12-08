#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log_info() {
  echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
  echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_error() {
  echo -e "${RED}[ERROR]${NC} $1"
}

# Detect OS
OS="$(uname -s)"
case "${OS}" in
Linux*) OS_TYPE=linux ;;
Darwin*) OS_TYPE=darwin ;;
*)
  log_error "Unsupported operating system: ${OS}"
  exit 1
  ;;
esac

# Detect Architecture
ARCH="$(uname -m)"
case "${ARCH}" in
x86_64) ARCH_TYPE=amd64 ;;
arm64) ARCH_TYPE=arm64 ;;
aarch64) ARCH_TYPE=arm64 ;;
*)
  log_error "Unsupported architecture: ${ARCH}"
  exit 1
  ;;
esac

log_info "Detected OS: ${OS_TYPE}, Architecture: ${ARCH_TYPE}"

# Version to install
VERSION="v0.1.2"

log_info "Installing version: ${VERSION}"

# Construct download URL
# Naming convention from .goreleaser.yaml: {{ .ProjectName }}-v{{ .Version }}-{{ .Os }}-{{ .Arch }}
VERSION_NUM="${VERSION#v}"
FILENAME="gitti-v${VERSION_NUM}-${OS_TYPE}-${ARCH_TYPE}.tar.gz"
DOWNLOAD_URL="https://github.com/gohyuhan/gitti/releases/download/${VERSION}/${FILENAME}"

log_info "Download URL: ${DOWNLOAD_URL}"

# Create temp directory
TMP_DIR=$(mktemp -d)
trap 'rm -rf "${TMP_DIR}"' EXIT

# Download
log_info "Downloading ${FILENAME}..."
curl -sL "${DOWNLOAD_URL}" -o "${TMP_DIR}/${FILENAME}"

# Extract
log_info "Extracting..."
tar -xzf "${TMP_DIR}/${FILENAME}" -C "${TMP_DIR}"

# Find binary
BINARY_PATH="${TMP_DIR}/gitti"
if [ ! -f "${BINARY_PATH}" ]; then
  # Try to find it if it was in a subdirectory or named differently?
  # Based on goreleaser, it should be at root of archive usually, or we can find it.
  BINARY_PATH=$(find "${TMP_DIR}" -type f -name "gitti" | head -n 1)
fi

if [ ! -f "${BINARY_PATH}" ]; then
  log_error "Binary 'gitti' not found in extracted archive."
  exit 1
fi

# Install
INSTALL_DIR="/usr/local/bin"
TARGET_PATH="${INSTALL_DIR}/gitti"

log_info "Installing to ${TARGET_PATH}..."

if [ -w "${INSTALL_DIR}" ]; then
  mv "${BINARY_PATH}" "${TARGET_PATH}"
else
  log_info "Requires sudo to install to ${INSTALL_DIR}"
  sudo mv "${BINARY_PATH}" "${TARGET_PATH}"
fi

sudo chmod +x "${TARGET_PATH}"

log_success "Gitti installed successfully!"
log_info "Run 'gitti --version' to verify."
