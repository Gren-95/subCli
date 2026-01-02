#!/bin/bash
# subCli installation script

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default installation directory
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="subCli"

# Check if running with sudo for system-wide install
if [ "$EUID" -ne 0 ] && [ "$1" != "--user" ]; then
    echo -e "${YELLOW}Note: Installing to system directory requires sudo${NC}"
    echo "Run with: sudo ./install.sh"
    echo "Or for user install: ./install.sh --user"
    exit 1
fi

# User installation option
if [ "$1" == "--user" ]; then
    INSTALL_DIR="$HOME/.local/bin"
    mkdir -p "$INSTALL_DIR"
    echo -e "${GREEN}Installing to user directory: $INSTALL_DIR${NC}"
else
    echo -e "${GREEN}Installing to system directory: $INSTALL_DIR${NC}"
fi

# Build the binary
echo "Building subCli..."
go build -o "$BINARY_NAME" || {
    echo -e "${RED}Build failed!${NC}"
    exit 1
}

# Check if binary exists and will be overwritten
if [ -f "$INSTALL_DIR/$BINARY_NAME" ]; then
    echo -e "${YELLOW}Existing installation found at $INSTALL_DIR/$BINARY_NAME${NC}"
    echo "Overwriting..."
fi

# Install the binary
echo "Installing $BINARY_NAME to $INSTALL_DIR..."
cp -f "$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME" || {
    echo -e "${RED}Installation failed!${NC}"
    exit 1
}

# Set executable permissions
chmod +x "$INSTALL_DIR/$BINARY_NAME"

# Verify installation
if [ -f "$INSTALL_DIR/$BINARY_NAME" ]; then
    echo -e "${GREEN}✓ Installation successful!${NC}"
    echo ""
    echo "Binary installed to: $INSTALL_DIR/$BINARY_NAME"
    
    # Check if directory is in PATH
    if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
        echo -e "${YELLOW}Warning: $INSTALL_DIR is not in your PATH${NC}"
        echo "Add it by running:"
        echo "  echo 'export PATH=\"\$PATH:$INSTALL_DIR\"' >> ~/.bashrc"
        echo "  source ~/.bashrc"
    fi
    
    echo ""
    echo "Run 'subCli setup' to configure your connection."
    echo "Then try: subCli --shuffle | mpv --playlist=-"
else
    echo -e "${RED}Installation verification failed!${NC}"
    exit 1
fi

