# Release Guide

## How to Create a Release

This project uses GitHub Actions with GoReleaser to automatically build and publish releases.

### Steps to Release

1. **Ensure everything is committed and pushed:**
   ```bash
   git add .
   git commit -m "Prepare for release"
   git push origin main
   ```

2. **Create and push a version tag:**
   ```bash
   # For a new version, e.g., v1.0.0
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

3. **GitHub Actions will automatically:**
   - Build binaries for Linux (amd64, arm64) and macOS (amd64, arm64)
   - Create a GitHub release
   - Upload the binaries and checksums
   - Generate release notes

### Version Numbering

Follow [Semantic Versioning](https://semver.org/):
- `v1.0.0` - Major release (breaking changes)
- `v1.1.0` - Minor release (new features)
- `v1.0.1` - Patch release (bug fixes)

### Example Release Workflow

```bash
# Make your changes
git add .
git commit -m "Add new feature"

# Create a new version tag
git tag -a v1.0.0 -m "Release v1.0.0 - Initial release"

# Push everything
git push origin main
git push origin v1.0.0
```

### Testing Releases Locally

Before creating a release, test the build locally:

```bash
# Install GoReleaser (if not already installed)
go install github.com/goreleaser/goreleaser@latest

# Test the release process (without publishing)
goreleaser release --snapshot --clean

# Check the dist/ directory for built artifacts
ls -lh dist/
```

### Manual Release (if needed)

If you need to create a release manually:

```bash
# Build for current platform
go build -o subcli

# Or build for specific platforms
GOOS=linux GOARCH=amd64 go build -o subcli-linux-amd64
GOOS=darwin GOARCH=amd64 go build -o subcli-darwin-amd64
```

## Release Checklist

Before creating a release:

- [ ] All tests pass
- [ ] Documentation is updated
- [ ] Version number follows semantic versioning
- [ ] CHANGELOG is updated (if you have one)
- [ ] All changes are committed and pushed

## Automated Workflows

### Release Workflow
Triggered on: Tag push (v*)
- Builds binaries for multiple platforms
- Creates GitHub release
- Uploads artifacts

### Build Workflow
Triggered on: Push to main/master, Pull requests
- Runs `go build`
- Runs `go test`
- Verifies binary works

## Download Released Binaries

Users can download binaries from:
```
https://github.com/gren-95/subCli/releases/latest
```

Or install using Go:
```bash
go install github.com/gren-95/subCli@latest
```

