#!/bin/bash
# Installation script for prom2sar

set -e

INSTALL_DIR="/usr/local/bin"
BINARY_NAME="prom2sar"

echo "=== prom2sar Installation ==="
echo ""

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    echo "This script requires sudo/root privileges to install to ${INSTALL_DIR}"
    echo "Please run: sudo ./install.sh"
    exit 1
fi

# Build the binary if it doesn't exist
if [ ! -f "bin/${BINARY_NAME}" ]; then
    echo "Binary not found. Building..."
    make build-cli
fi

# Verify binary exists
if [ ! -f "bin/${BINARY_NAME}" ]; then
    echo "Error: Failed to build binary"
    exit 1
fi

# Install
echo "Installing ${BINARY_NAME} to ${INSTALL_DIR}..."
cp "bin/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
chmod +x "${INSTALL_DIR}/${BINARY_NAME}"

# Verify installation
if command -v ${BINARY_NAME} &> /dev/null; then
    echo ""
    echo "✓ Installation successful!"
    echo ""
    echo "Installed: $(which ${BINARY_NAME})"
    echo "Version:   $(${BINARY_NAME} --version)"
    echo ""
    echo "Try it:"
    echo "  ${BINARY_NAME} --help"
    echo ""
else
    echo "Error: Installation failed"
    exit 1
fi
