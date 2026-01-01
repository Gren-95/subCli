#!/bin/bash
# Example script for using subCli with mpv

# Display usage information
usage() {
    cat << EOF
Usage: $0 [OPTIONS]

Examples:
    $0                          # Play random albums
    $0 --shuffle                # Play random albums shuffled
    $0 --favorites              # Play your favorites
    $0 --search "rock" --type song --shuffle  # Search and shuffle

Options are passed directly to subCli. Run 'subcli --help' for all options.
EOF
    exit 1
}

# Check if subcli exists
if ! command -v subcli &> /dev/null; then
    echo "Error: subcli not found in PATH"
    echo "Please build it with: go build -o subcli"
    exit 1
fi

# Check if mpv exists
if ! command -v mpv &> /dev/null; then
    echo "Error: mpv not found in PATH"
    echo "Please install mpv"
    exit 1
fi

# Show help if requested
if [[ "$1" == "-h" ]] || [[ "$1" == "--help" ]]; then
    usage
fi

# Run subcli and pipe to mpv
echo "Starting music playback..."
echo "Commands: subcli $@"
echo "Press Ctrl+C to stop"
echo ""

subcli "$@" | mpv --playlist=- --no-video --term-osd-bar

