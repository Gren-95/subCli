package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gren-95/subCli/internal/api"
	"github.com/gren-95/subCli/internal/config"
	"github.com/spf13/cobra"
)

var (
	shuffle    bool
	playlist   string
	album      string
	artist     string
	search     string
	searchType string
	limit      int
	favorites  bool
	m3uFormat  bool
)

func main() {
	rand.Seed(time.Now().UnixNano())

	rootCmd := &cobra.Command{
		Use:   "subCli",
		Short: "A CLI for streaming music from Subsonic-compatible servers to mpv",
		Long:  `subCli allows you to search, queue, and stream music from Subsonic servers directly to mpv.`,
		Run:   runCLI,
	}

	// Setup command for first-time configuration
	setupCmd := &cobra.Command{
		Use:   "setup",
		Short: "Run interactive setup to configure subCli",
		Long:  `Interactively configure subCli with your Subsonic server credentials.`,
		Run:   runSetup,
	}

	rootCmd.AddCommand(setupCmd)

	rootCmd.Flags().BoolVarP(&shuffle, "shuffle", "s", false, "Shuffle the playlist")
	rootCmd.Flags().StringVarP(&playlist, "playlist", "p", "", "Play a specific playlist by name or ID")
	rootCmd.Flags().StringVarP(&album, "album", "a", "", "Play a specific album by ID")
	rootCmd.Flags().StringVarP(&artist, "artist", "r", "", "Play albums from a specific artist by ID")
	rootCmd.Flags().StringVarP(&search, "search", "q", "", "Search for songs/albums/artists")
	rootCmd.Flags().StringVarP(&searchType, "type", "t", "song", "Search type: song, album, artist")
	rootCmd.Flags().IntVarP(&limit, "limit", "n", 50, "Limit number of results")
	rootCmd.Flags().BoolVarP(&favorites, "favorites", "f", false, "Play favorite songs")
	rootCmd.Flags().BoolVarP(&m3uFormat, "m3u", "m", false, "Output in M3U playlist format with metadata")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runSetup(cmd *cobra.Command, args []string) {
	if err := config.InteractiveSetup(); err != nil {
		fmt.Fprintf(os.Stderr, "Setup error: %v\n", err)
		os.Exit(1)
	}

	// Test the connection
	if err := api.SubsonicPing(); err != nil {
		fmt.Fprintf(os.Stderr, "✗ Connection test failed: %v\n", err)
		fmt.Fprintf(os.Stderr, "Please check your credentials and try again.\n")
		os.Exit(1)
	}

	fmt.Println("✓ Connection test successful!")
	fmt.Println()
	fmt.Println("You're all set! Try running:")
	fmt.Println("  subCli --shuffle | mpv --playlist=-")
}

func runCLI(cmd *cobra.Command, args []string) {
	// Load configuration
	if err := api.LoadConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		fmt.Fprintf(os.Stderr, "\nRun 'subCli setup' to configure your connection.\n")
		os.Exit(1)
	}

	// Validate configuration
	if api.AppConfig.Username == "" || api.AppConfig.URL == "" {
		fmt.Fprintf(os.Stderr, "Error: Configuration is incomplete\n")
		fmt.Fprintf(os.Stderr, "Run 'subCli setup' to configure your connection.\n")
		os.Exit(1)
	}

	if config.GetPassword() == "" {
		fmt.Fprintf(os.Stderr, "Error: Password not configured\n")
		fmt.Fprintf(os.Stderr, "Run 'subCli setup' to configure your connection.\n")
		os.Exit(1)
	}

	// Test connection
	if err := api.SubsonicPing(); err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to server: %v\n", err)
		os.Exit(1)
	}

	var songs []api.Song
	var err error

	// Determine what to play based on flags
	if favorites {
		songs, err = getFavorites()
	} else if playlist != "" {
		songs, err = getPlaylistSongs(playlist)
	} else if album != "" {
		songs, err = api.SubsonicGetAlbum(album)
	} else if artist != "" {
		songs, err = getArtistSongs(artist)
	} else if search != "" {
		songs, err = searchAndGetSongs(search, searchType)
	} else {
		// Default: get random albums
		songs, err = getRandomAlbums()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching songs: %v\n", err)
		os.Exit(1)
	}

	if len(songs) == 0 {
		fmt.Fprintf(os.Stderr, "No songs found\n")
		os.Exit(1)
	}

	// Apply limit if specified
	if limit > 0 && len(songs) > limit {
		songs = songs[:limit]
	}

	// Shuffle if requested
	if shuffle {
		rand.Shuffle(len(songs), func(i, j int) {
			songs[i], songs[j] = songs[j], songs[i]
		})
	}

	// Output URLs for mpv
	if m3uFormat {
		outputM3U(songs)
	} else {
		outputURLs(songs)
	}
}

func getFavorites() ([]api.Song, error) {
	result, err := api.SubsonicGetStarred()
	if err != nil {
		return nil, err
	}
	return result.Songs, nil
}

func getPlaylistSongs(playlistNameOrID string) ([]api.Song, error) {
	// First, try to use it as an ID
	songs, err := api.SubsonicGetPlaylistSongs(playlistNameOrID)
	if err == nil {
		return songs, nil
	}

	// If that fails, search by name
	playlists, err := api.SubsonicGetPlaylists()
	if err != nil {
		return nil, err
	}

	for _, pl := range playlists {
		if strings.EqualFold(pl.Name, playlistNameOrID) || pl.ID == playlistNameOrID {
			return api.SubsonicGetPlaylistSongs(pl.ID)
		}
	}

	return nil, fmt.Errorf("playlist not found: %s", playlistNameOrID)
}

func getArtistSongs(artistID string) ([]api.Song, error) {
	albums, err := api.SubsonicGetArtist(artistID)
	if err != nil {
		return nil, err
	}

	var allSongs []api.Song
	for _, album := range albums {
		songs, err := api.SubsonicGetAlbum(album.ID)
		if err != nil {
			continue
		}
		allSongs = append(allSongs, songs...)
	}

	return allSongs, nil
}

func searchAndGetSongs(query string, searchType string) ([]api.Song, error) {
	switch strings.ToLower(searchType) {
	case "song":
		return api.SubsonicSearchSong(query, 0)
	case "album":
		albums, err := api.SubsonicSearchAlbum(query, 0)
		if err != nil {
			return nil, err
		}
		if len(albums) == 0 {
			return nil, nil
		}
		// Get songs from the first matching album
		return api.SubsonicGetAlbum(albums[0].ID)
	case "artist":
		artists, err := api.SubsonicSearchArtist(query, 0)
		if err != nil {
			return nil, err
		}
		if len(artists) == 0 {
			return nil, nil
		}
		// Get songs from the first matching artist
		return getArtistSongs(artists[0].ID)
	default:
		return nil, fmt.Errorf("invalid search type: %s (use song, album, or artist)", searchType)
	}
}

func getRandomAlbums() ([]api.Song, error) {
	// Fetch random albums progressively to create a buffer
	// This gets albums in batches for better performance
	batchCount := 10 // Fetch 10 random albums at a time
	
	albums, err := api.SubsonicGetAlbumList("random")
	if err != nil {
		return nil, err
	}

	// Limit to first batch for buffering approach
	if len(albums) > batchCount {
		albums = albums[:batchCount]
	}

	var allSongs []api.Song
	seenSongs := make(map[string]bool) // Track songs to avoid duplicates
	
	for _, album := range albums {
		songs, err := api.SubsonicGetAlbum(album.ID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to get album %s: %v\n", album.Name, err)
			continue
		}
		
		// Add songs, filtering out duplicates
		for _, song := range songs {
			if !seenSongs[song.ID] {
				allSongs = append(allSongs, song)
				seenSongs[song.ID] = true
			}
		}
	}

	return allSongs, nil
}

func outputURLs(songs []api.Song) {
	// Output URLs one at a time as they're generated
	for _, song := range songs {
		url := api.SubsonicStream(song.ID)
		fmt.Println(url)
	}
}

func outputM3U(songs []api.Song) {
	// Output in M3U format with metadata
	fmt.Println("#EXTM3U")
	for _, song := range songs {
		// #EXTINF:duration,Artist - Title
		fmt.Printf("#EXTINF:%d,%s - %s\n", song.Duration, song.Artist, song.Title)
		url := api.SubsonicStream(song.ID)
		fmt.Println(url)
	}
}
