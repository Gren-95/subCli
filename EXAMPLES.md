# subCli Usage Examples

## Basic Playback

### Play random albums
```bash
subcli | mpv --playlist=-
```

### Play random albums (shuffled)
```bash
subcli --shuffle | mpv --playlist=-
```

### Play with limited number of songs
```bash
subcli --limit 20 | mpv --playlist=-
```

## Search Examples

### Search for songs
```bash
subcli --search "bohemian rhapsody" --type song | mpv --playlist=-
```

### Search for albums and shuffle
```bash
subcli --search "dark side of the moon" --type album --shuffle | mpv --playlist=-
```

### Search for artist's music
```bash
subcli --search "miles davis" --type artist --shuffle | mpv --playlist=-
```

### Search with limit
```bash
subcli --search "jazz" --type song --limit 30 --shuffle | mpv --playlist=-
```

## Playlists

### Play a playlist by name
```bash
subcli --playlist "Chill Mix" | mpv --playlist=-
```

### Play a playlist by ID (shuffled)
```bash
subcli --playlist "pl-12345" --shuffle | mpv --playlist=-
```

## Albums and Artists

### Play a specific album
```bash
subcli --album "al-67890" | mpv --playlist=-
```

### Play all songs from an artist
```bash
subcli --artist "ar-11111" --shuffle | mpv --playlist=-
```

## Favorites

### Play your favorite songs
```bash
subcli --favorites | mpv --playlist=-
```

### Shuffle favorites with limit
```bash
subcli --favorites --shuffle --limit 50 | mpv --playlist=-
```

## MPV Options

### Play without video window
```bash
subcli --shuffle | mpv --playlist=- --no-video
```

### Play with volume control
```bash
subcli --favorites | mpv --playlist=- --volume=50
```

### Play in background
```bash
subcli --shuffle | mpv --playlist=- --no-video &
```

### Show terminal OSD bar
```bash
subcli --shuffle | mpv --playlist=- --term-osd-bar
```

### Play with MPV socket for control
```bash
subcli --shuffle | mpv --playlist=- --input-ipc-server=/tmp/mpvsocket
```

Then control it:
```bash
# Pause/Resume
echo '{"command": ["cycle", "pause"]}' | socat - /tmp/mpvsocket

# Next track
echo '{"command": ["playlist-next"]}' | socat - /tmp/mpvsocket

# Get current position
echo '{"command": ["get_property", "time-pos"]}' | socat - /tmp/mpvsocket
```

## Saving Playlists

### Save as plain URL list
```bash
subcli --favorites > favorites.txt
```

### Save as M3U playlist with metadata
```bash
subcli --favorites --m3u > favorites.m3u
```

### Save shuffled playlist for later
```bash
subcli --search "rock" --type song --shuffle --limit 100 --m3u > rock_playlist.m3u
mpv --playlist=rock_playlist.m3u
```

## Advanced Examples

### Daily random music routine
```bash
#!/bin/bash
# Add this to your startup scripts
subcli --shuffle --limit 100 | mpv --playlist=- --no-video --volume=30 &
```

### Create genre-based playlists
```bash
# Create multiple playlists
subcli --search "jazz" --type song --limit 50 --m3u > jazz.m3u
subcli --search "rock" --type song --limit 50 --m3u > rock.m3u
subcli --search "classical" --type song --limit 50 --m3u > classical.m3u

# Play them in sequence
mpv --playlist=jazz.m3u --playlist=rock.m3u --playlist=classical.m3u
```

### Random artist discovery
```bash
# Get random albums and extract unique artists
subcli --shuffle --limit 200 | sort -u | head -50 | mpv --playlist=-
```

### Quiet background music
```bash
subcli --search "ambient" --type song --shuffle | \
    mpv --playlist=- --no-video --volume=20 --loop-playlist &
```

### Work/Study session
```bash
# Play 2 hours of music
subcli --shuffle --limit 40 | mpv --playlist=- --no-video --term-osd-bar
```

## Using the Helper Script

The `play.sh` script simplifies common usage:

```bash
# Random music
./play.sh --shuffle

# Favorites
./play.sh --favorites --shuffle

# Search
./play.sh --search "chillhop" --type song --shuffle

# Custom MPV options can be added to the script
```

## Combining with Other Tools

### With fzf for interactive selection
```bash
# This is conceptual - subcli outputs URLs, not song names
# You'd need to modify it to output JSON or parse the M3U format
subcli --favorites --m3u | grep -v "^#" | mpv --playlist=-
```

### With notification on track change
```bash
# Using mpv's title property (mpv handles this with --term-osd-bar)
subcli --shuffle | mpv --playlist=- --term-osd-bar
```

### With scrobbling (if supported by your Subsonic server)
Subsonic API handles scrobbling automatically when you stream!

## Tips

1. **Use --limit** to avoid overwhelming playlists
2. **Use --m3u** to save playlists with metadata for better organization
3. **Use --shuffle** for variety in long listening sessions
4. **Combine with mpv's --loop-playlist** for endless playback
5. **Use mpv's socket** for runtime control without restarting

## Troubleshooting

If subcli outputs URLs but mpv doesn't play:
```bash
# Test a single URL
subcli --favorites --limit 1 | xargs mpv

# Check if URLs are valid
subcli --favorites --limit 1

# Test with verbose mpv output
subcli --shuffle --limit 5 | mpv --playlist=- --msg-level=all=debug
```

