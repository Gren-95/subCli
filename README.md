# subCli

A command-line interface for streaming music from Subsonic-compatible servers to mpv.

## Features

- 🎵 Stream music directly to mpv via pipe
- 🔀 Shuffle playlists
- 🔁 Loop modes (none, all, one)
- ⭐ Access favorites
- 🔍 Search for songs, albums, and artists
- 📋 Play playlists
- 🎲 Random album playback

## Installation

### From Releases (Recommended)

Download the latest binary for your platform:

```bash
# Go to releases page
https://github.com/gren-95/subCli/releases/latest
```

### Using Go

```bash
go install github.com/gren-95/subCli@latest
```

### From Source

```bash
git clone https://github.com/gren-95/subCli
cd subCli
go build -o subcli
```

> **Note:** Releases are automatically created on every push to main using semantic versioning based on commit messages.

## Configuration

Create a configuration file at `~/.config/subcli/config.yaml`:

```yaml
username: your_username
password: your_password
URL: https://your-subsonic-server.com
```

## Usage

### Basic Usage

Pipe to mpv or VLC to play music:

```bash
# With mpv
subcli | mpv --playlist=-

# With VLC
subcli --m3u | vlc -
```

### Search and Play

Search for songs and play them:

```bash
subcli --search "artist name" --type song | mpv --playlist=-
```

Search for an album:

```bash
subcli --search "album name" --type album | mpv --playlist=-
```

Search for an artist:

```bash
subcli --search "artist name" --type artist | mpv --playlist=-
```

### Shuffle and Loop

Shuffle your playlist:

```bash
subcli --search "rock" --shuffle | mpv --playlist=-
```

Loop all songs:

```bash
subcli --playlist "My Favorites" --loop all | mpv --playlist=-
```

Loop one song:

```bash
subcli --search "favorite song" --loop one | mpv --playlist=-
```

### Play Playlists

Play a playlist by name:

```bash
subcli --playlist "Chill Mix" | mpv --playlist=-
```

Play a playlist by ID:

```bash
subcli --playlist "playlist-id-123" | mpv --playlist=-
```

### Play Albums

Play a specific album by ID:

```bash
subcli --album "album-id-456" | mpv --playlist=-
```

### Play Artist's Music

Play all songs from an artist by ID:

```bash
subcli --artist "artist-id-789" | mpv --playlist=-
```

### Play Favorites

Play your starred/favorite songs:

```bash
subcli --favorites | mpv --playlist=-
```

Shuffle your favorites:

```bash
subcli --favorites --shuffle | mpv --playlist=-
```

### Limit Results

Limit the number of songs:

```bash
subcli --search "pop" --limit 20 | mpv --playlist=-
```

### Random Albums

If you don't specify any flags, subCli will fetch random albums:

```bash
subcli | mpv --playlist=-
```

Shuffle random albums:

```bash
subcli --shuffle | mpv --playlist=-
```

## Command-Line Options

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--search` | `-q` | Search for songs/albums/artists | - |
| `--type` | `-t` | Search type: song, album, artist | song |
| `--playlist` | `-p` | Play a specific playlist by name or ID | - |
| `--album` | `-a` | Play a specific album by ID | - |
| `--artist` | `-r` | Play albums from a specific artist by ID | - |
| `--favorites` | `-f` | Play favorite songs | false |
| `--shuffle` | `-s` | Shuffle the playlist | false |
| `--loop` | `-l` | Loop mode: none, all, one | none |
| `--limit` | `-n` | Limit number of results | 50 |
| `--m3u` | `-m` | Output in M3U playlist format | false |

## Media Player Tips

### MPV

**Background Playback:**

```bash
subcli --shuffle | mpv --playlist=- &
```

**MPV Socket Control:**

For better control, use mpv with a socket:

```bash
subcli --shuffle | mpv --playlist=- --input-ipc-server=/tmp/mpvsocket
```

Then control playback with:

```bash
echo '{ "command": ["cycle", "pause"] }' | socat - /tmp/mpvsocket
```

### VLC

**Basic VLC usage:**

```bash
# With M3U format (recommended)
subcli --shuffle --m3u | vlc -

# Background playback
subcli --shuffle --m3u | vlc --intf dummy -

# Save and play
subcli --favorites --m3u > playlist.m3u
vlc playlist.m3u
```

**VLC HTTP Interface:**

```bash
# Start with web interface
subcli --shuffle --m3u | vlc --intf http --http-password mypassword -

# Access at http://localhost:8080
```

### Save Playlist

Save the playlist to a file:

```bash
subcli --favorites > playlist.m3u
```

Then play it later:

```bash
mpv --playlist=playlist.m3u
```

## Examples

Random shuffled music for background listening:

```bash
subcli --shuffle --limit 100 | mpv --playlist=- --no-video --volume=50
```

Play your favorites on loop:

```bash
subcli --favorites --shuffle --loop all | mpv --playlist=-
```

Quick search and play:

```bash
subcli -q "miles davis" -t artist -s | mpv --playlist=-
```

## License

MIT License - see LICENSE file for details

## Credits

Originally based on [SubTUI](https://github.com/MattiaPun/SubTUI), converted to a CLI application by [gren-95](https://github.com/gren-95).

