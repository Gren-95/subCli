# Quick Start Guide for subCli

## First Time Setup

1. **Create configuration file:**

```bash
mkdir -p ~/.config/subcli
cat > ~/.config/subcli/config.yaml << EOF
username: your_username
password: your_password
URL: https://your-subsonic-server.com
EOF
```

2. **Build the application:**

```bash
go build -o subcli
```

3. **Test connection:**

```bash
./subcli --favorites --limit 1
```

This should output a single URL if everything is configured correctly.

## Basic Usage Examples

### Play random music (shuffled)

```bash
./subcli --shuffle | mpv --playlist=-
```

Or use the helper script:

```bash
./play.sh --shuffle
```

### Search and play

```bash
./subcli --search "pink floyd" --type artist --shuffle | mpv --playlist=-
```

### Play a playlist

```bash
./subcli --playlist "Road Trip" | mpv --playlist=-
```

### Play favorites (shuffled)

```bash
./subcli --favorites --shuffle | mpv --playlist=-
```

### Limit songs and shuffle

```bash
./subcli --search "jazz" --limit 30 --shuffle | mpv --playlist=-
```

## Understanding the Output

subCli outputs stream URLs, one per line:

```
https://music.example.com/rest/stream?id=123&...
https://music.example.com/rest/stream?id=456&...
https://music.example.com/rest/stream?id=789&...
```

These URLs can be:
- Piped directly to mpv: `subcli ... | mpv --playlist=-`
- Saved to a file: `subcli ... > playlist.m3u`
- Used with other media players that support URL playlists

## Tips

### Background Playback

Run mpv in the background:

```bash
./subcli --shuffle | mpv --playlist=- --no-video &
```

### Save Playlist for Later

```bash
./subcli --favorites --shuffle > ~/Music/favorites.m3u
mpv --playlist=~/Music/favorites.m3u
```

### Control MPV While Playing

Use mpv's IPC socket:

```bash
./subcli --shuffle | mpv --playlist=- --input-ipc-server=/tmp/mpvsocket
```

Then in another terminal:

```bash
# Pause/Resume
echo '{"command": ["cycle", "pause"]}' | socat - /tmp/mpvsocket

# Next track
echo '{"command": ["playlist-next"]}' | socat - /tmp/mpvsocket

# Previous track
echo '{"command": ["playlist-prev"]}' | socat - /tmp/mpvsocket
```

## Troubleshooting

### Configuration Error

If you see "Error loading config", make sure:
- The config file exists at `~/.config/subcli/config.yaml`
- The YAML syntax is correct (watch for spaces, not tabs)
- All three fields are filled: username, password, URL

### Connection Error

If you see "Error connecting to server":
- Check your URL is correct (include http:// or https://)
- Verify your username and password
- Make sure your Subsonic server is running and accessible

### No Songs Found

If you see "No songs found":
- Try a different search query
- Check if the playlist/album/artist ID is correct
- Try `--favorites` to see if you can access any content

### MPV Not Playing

If URLs are printed but mpv doesn't play:
- Make sure mpv is installed: `which mpv`
- Try playing one URL directly: `mpv "URL_HERE"`
- Check your network connection to the Subsonic server

