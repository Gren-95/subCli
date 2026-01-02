#!/bin/bash
# subCli uninstallation script

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

BINARY_NAME="subCli"

# Check both system and user installation locations
SYSTEM_PATH="/usr/local/bin/$BINARY_NAME"
USER_PATH="$HOME/.local/bin/$BINARY_NAME"

removed=false

# Remove from system location
if [ -f "$SYSTEM_PATH" ]; then
    if [ "$EUID" -ne 0 ]; then
        echo -e "${YELLOW}System installation found. Need sudo to remove.${NC}"
        echo "Run: sudo ./uninstall.sh"
        exit 1
    fi
    
    echo "Removing $SYSTEM_PATH..."
    rm -f "$SYSTEM_PATH"
    echo -e "${GREEN}✓ Removed system installation${NC}"
    removed=true
fi

# Remove from user location
if [ -f "$USER_PATH" ]; then
    echo "Removing $USER_PATH..."
    rm -f "$USER_PATH"
    echo -e "${GREEN}✓ Removed user installation${NC}"
    removed=true
fi

if [ "$removed" = false ]; then
    echo -e "${YELLOW}No installation found${NC}"
else
    echo ""
    echo "Configuration file preserved at: ~/.config/subCli/config.yaml"
    echo "To remove config: rm -rf ~/.config/subCli"
fi

