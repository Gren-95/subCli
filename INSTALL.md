# Installation Guide

## Installation Methods

### 1. Using Install Script (Recommended)

#### System-wide Installation
Install to `/usr/local/bin` (requires sudo):

```bash
git clone https://github.com/gren-95/subCli
cd subCli
sudo ./install.sh
```

This will:
- Build the binary
- Install to `/usr/local/bin/subCli`
- Overwrite existing installation if present
- Set proper executable permissions

#### User Installation
Install to `~/.local/bin` (no sudo required):

```bash
git clone https://github.com/gren-95/subCli
cd subCli
./install.sh --user
```

### 2. Manual Installation

```bash
# Build
go build -o subCli

# System-wide
sudo cp subCli /usr/local/bin/
sudo chmod +x /usr/local/bin/subCli

# User-only
mkdir -p ~/.local/bin
cp subCli ~/.local/bin/
chmod +x ~/.local/bin/subCli
```

### 3. Using Go Install

```bash
go install github.com/gren-95/subCli@latest
```

Binary will be installed to `$GOPATH/bin` or `~/go/bin`.

### 4. From Releases

Download pre-built binaries from:
```
https://github.com/gren-95/subCli/releases/latest
```

Then:
```bash
# Extract and install
tar -xzf subCli_Linux_x86_64.tar.gz
sudo cp subCli /usr/local/bin/
sudo chmod +x /usr/local/bin/subCli
```

## Post-Installation

### Verify Installation

```bash
which subCli
subCli --help
```

### Configure

Run the setup wizard:

```bash
subCli setup
```

### PATH Configuration

If `subCli` is not found after installation:

**For `/usr/local/bin` (usually already in PATH):**
```bash
echo $PATH | grep /usr/local/bin
```

**For `~/.local/bin`:**
Add to your shell config (`~/.bashrc`, `~/.zshrc`, etc.):
```bash
export PATH="$HOME/.local/bin:$PATH"
```

Then reload:
```bash
source ~/.bashrc  # or ~/.zshrc
```

## Uninstallation

Use the uninstall script:

```bash
# System installation
sudo ./uninstall.sh

# User installation
./uninstall.sh
```

Or manually:

```bash
# Remove binary
sudo rm /usr/local/bin/subCli
# or
rm ~/.local/bin/subCli

# Optionally remove config
rm -rf ~/.config/subCli
```

## Upgrading

To upgrade an existing installation:

```bash
cd subCli
git pull origin main
sudo ./install.sh  # or ./install.sh --user
```

The install script will automatically overwrite the existing binary.

## Troubleshooting

### Permission Denied

If you see "Permission denied" when running `subCli`:
```bash
chmod +x /usr/local/bin/subCli
# or
chmod +x ~/.local/bin/subCli
```

### Command Not Found

1. Check if installed:
   ```bash
   ls -l /usr/local/bin/subCli
   ls -l ~/.local/bin/subCli
   ```

2. Check PATH:
   ```bash
   echo $PATH
   ```

3. Add to PATH if needed (see PATH Configuration above)

### Install Script Fails

Make sure you have:
- Go installed (`go version`)
- Git installed (`git --version`)
- Proper permissions (use `sudo` for system install)

### Multiple Installations

If you have both system and user installations, the one in your PATH first will be used.
Check with:
```bash
which subCli
```

To remove all installations:
```bash
sudo rm /usr/local/bin/subCli
rm ~/.local/bin/subCli
```

