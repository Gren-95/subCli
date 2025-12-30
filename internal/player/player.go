package player

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"github.com/MattiaPun/SubTUI/internal/api"
	"github.com/gdrens/mpv"
)

var (
	mpvClient *mpv.Client
	mpvCmd    *exec.Cmd
)

type PlayerStatus struct {
	Title    string
	Artist   string
	Album    string
	Current  float64
	Duration float64
	Paused   bool
	Volume   float64
}

func InitPlayer() error {
	socketPath := filepath.Join(os.TempDir(), fmt.Sprintf("subtui_mpv_socket_%d", os.Getuid()))

	_ = exec.Command("pkill", "-f", socketPath).Run()
	time.Sleep(200 * time.Millisecond)

	args := []string{
		"--idle",
		"--no-video",
		"--input-ipc-server=" + socketPath,
	}

	mpvCmd = exec.Command("mpv", args...)
	if err := mpvCmd.Start(); err != nil {
		return fmt.Errorf("failed to start mpv: %v", err)
	}

	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		if _, err := os.Stat(socketPath); err == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	ipcc := mpv.NewIPCClient(socketPath)
	client := mpv.NewClient(ipcc)
	mpvClient = client

	return nil
}

func ShutdownPlayer() {
	if mpvCmd != nil {
		_ = mpvCmd.Process.Signal(syscall.SIGTERM)
	}
}

func PlaySong(songID string) error {
	if mpvClient == nil {
		return fmt.Errorf("player not initialized")
	}

	url := api.SubsonicStream(songID)
	if err := mpvClient.LoadFile(url, mpv.LoadFileModeReplace); err != nil {
		return err
	}

	api.SubsonicScrobble(songID, false)

	_ = mpvClient.SetProperty("pause", false)

	return nil
}

func TogglePause() {
	if mpvClient == nil {
		return
	}

	status := mpvClient.IsPause()
	_ = mpvClient.SetProperty("pause", !status)
}

func Back10Seconds() {
	_ = mpvClient.Seek(-10)
}

func Forward10Seconds() {
	_ = mpvClient.Seek(+10)
}

func GetPlayerStatus() PlayerStatus {
	if mpvClient == nil {
		return PlayerStatus{}
	}

	title := mpvClient.GetProperty("media-title")
	artist := mpvClient.GetProperty("metadata/by-key/artist")
	album := mpvClient.GetProperty("metadata/by-key/album")

	pos := mpvClient.Position()
	dur := mpvClient.Duration()
	paused := mpvClient.IsPause()
	vol, _ := mpvClient.GetFloatProperty("volume")

	return PlayerStatus{
		Title:    fmt.Sprintf("%v", title),
		Artist:   fmt.Sprintf("%v", artist),
		Album:    fmt.Sprintf("%v", album),
		Current:  pos,
		Duration: dur,
		Paused:   paused,
		Volume:   vol,
	}
}
