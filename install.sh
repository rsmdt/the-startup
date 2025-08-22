#!/usr/bin/env bash
# The Startup - Installer Wrapper
# Downloads and runs the Go binary installer
#
# Usage:
#   curl -LsSf https://raw.githubusercontent.com/rsmdt/the-startup/main/install.sh | sh
#   curl -LsSf https://raw.githubusercontent.com/rsmdt/the-startup/main/install.sh | sh -s -- --yes
#   curl -LsSf https://raw.githubusercontent.com/rsmdt/the-startup/main/install.sh | sh -s -- --local
#   curl -LsSf https://raw.githubusercontent.com/rsmdt/the-startup/main/install.sh | sh -s -- -ly
#
# All flags are passed through to the binary installer.

set -e

# Configuration
REPO_OWNER="rsmdt"
REPO_NAME="the-startup"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Map architecture names
case "$ARCH" in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo "Error: Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

# Construct download URL
BINARY_NAME="${REPO_NAME}-${OS}-${ARCH}"
if [ "$OS" = "windows" ]; then
    BINARY_NAME="${BINARY_NAME}.exe"
fi

DOWNLOAD_URL="https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/latest/download/${BINARY_NAME}"

# Download the binary
echo "Downloading The Startup installer..."
echo "Platform: ${OS}-${ARCH}"

# Create temp file
TEMP_FILE=$(mktemp)
trap "rm -f $TEMP_FILE" EXIT

# Download with curl
if ! curl -fsSL "$DOWNLOAD_URL" -o "$TEMP_FILE"; then
    echo "Error: Failed to download installer from $DOWNLOAD_URL"
    echo "Please check your internet connection and try again."
    exit 1
fi

# Make executable
chmod +x "$TEMP_FILE"

# Run the installer
if [ $# -gt 0 ]; then
    echo "Starting installation with flags: $*"
else
    echo "Starting interactive installation..."
fi
"$TEMP_FILE" install "$@"

# Cleanup is handled by trap