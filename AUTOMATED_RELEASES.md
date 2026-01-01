# Automated Releases

This repository uses automated releases based on commit messages.

## How It Works

Every push to `main` automatically:
1. Analyzes your commit messages
2. Determines the version bump type
3. Creates a new tag and release
4. Builds binaries for multiple platforms

## Commit Message Format

Use [Conventional Commits](https://www.conventionalcommits.org/) to control version bumping:

### Patch Release (v1.0.0 → v1.0.1)
Bug fixes and minor changes:
```bash
git commit -m "fix: correct stream URL generation"
git commit -m "docs: update README"
git commit -m "chore: update dependencies"
```

### Minor Release (v1.0.0 → v1.1.0)
New features (backwards compatible):
```bash
git commit -m "feat: add support for VLC player"
git commit -m "feat(cli): add batch-size flag"
```

### Major Release (v1.0.0 → v2.0.0)
Breaking changes:
```bash
git commit -m "feat!: change API structure"
git commit -m "fix!: remove deprecated flags"
```

Or with body:
```bash
git commit -m "feat: new feature

BREAKING CHANGE: This changes the config file format"
```

## Examples

### Adding a Feature
```bash
git add .
git commit -m "feat: add playlist shuffle mode"
git push origin main
```
→ Creates v1.1.0 automatically

### Fixing a Bug
```bash
git add .
git commit -m "fix: resolve connection timeout issue"
git push origin main
```
→ Creates v1.0.1 automatically

### Multiple Commits
```bash
git commit -m "docs: improve examples"
git commit -m "fix: handle empty playlists"
git commit -m "feat: add random album selection"
git push origin main
```
→ Creates v1.1.0 (highest bump type wins)

## Commit Types

| Type | Description | Version Bump |
|------|-------------|--------------|
| `feat` | New feature | Minor (v1.1.0) |
| `fix` | Bug fix | Patch (v1.0.1) |
| `docs` | Documentation only | Patch (v1.0.1) |
| `style` | Code style/formatting | Patch (v1.0.1) |
| `refactor` | Code refactoring | Patch (v1.0.1) |
| `perf` | Performance improvement | Patch (v1.0.1) |
| `test` | Adding tests | Patch (v1.0.1) |
| `chore` | Maintenance | Patch (v1.0.1) |
| `!` or `BREAKING CHANGE:` | Breaking change | Major (v2.0.0) |

## Workflow Summary

1. **You commit and push to main**
2. GitHub Actions automatically:
   - Analyzes commit messages
   - Bumps version appropriately
   - Creates and pushes a git tag
   - Runs GoReleaser
   - Publishes binaries to GitHub Releases

## No Manual Tagging Needed!

❌ Don't do this:
```bash
git tag v1.0.0
git push origin v1.0.0
```

✅ Just push to main:
```bash
git push origin main
```

The version is automatically determined from your commits!

## Viewing Releases

All releases are automatically published at:
```
https://github.com/gren-95/subCli/releases
```

## First Release

The first time the workflow runs, it will:
- Start from v0.0.0
- Create v0.0.1 (or higher based on commits)
- Subsequent pushes will increment from there

## Disabling Auto-Release

If you need to push without creating a release, add `[skip ci]` to your commit:
```bash
git commit -m "docs: minor typo fix [skip ci]"
```

