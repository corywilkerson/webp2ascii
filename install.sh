#!/bin/bash
set -e

REPO="yourusername/webp2ascii"
VERSION="latest"
INSTALL_DIR="/usr/local/bin"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Download URL
if [ "$VERSION" = "latest" ]; then
    URL="https://github.com/${REPO}/releases/latest/download/webp2ascii-1.0.0-${OS}-${ARCH}.tar.gz"
else
    URL="https://github.com/${REPO}/releases/download/${VERSION}/webp2ascii-1.0.0-${OS}-${ARCH}.tar.gz"
fi

echo "Downloading webp2ascii..."
curl -sSL "$URL" | tar -xz -C /tmp/

echo "Installing to ${INSTALL_DIR}..."
sudo mv /tmp/webp2ascii-1.0.0-${OS}-${ARCH} "$INSTALL_DIR/webp2ascii"
sudo chmod +x "${INSTALL_DIR}/webp2ascii"

echo "✅ webp2ascii installed successfully!"
echo "Run 'webp2ascii -h' to get started"