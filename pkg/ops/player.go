package ops

import (
	"fmt"
	"strconv"

	"github.com/jrogala/heos-cli/client"
)

// PlayerInfo holds player details.
type PlayerInfo struct {
	Name    string `json:"name"`
	PID     string `json:"pid"`
	Model   string `json:"model"`
	Version string `json:"version"`
	Network string `json:"network"`
	Serial  string `json:"serial,omitempty"`
}

// NowPlaying holds currently playing media.
type NowPlaying struct {
	Type    string `json:"type"`
	Song    string `json:"song"`
	Artist  string `json:"artist"`
	Album   string `json:"album"`
	Station string `json:"station,omitempty"`
}

// PlayMode holds repeat and shuffle settings.
type PlayMode struct {
	Repeat  string `json:"repeat"`
	Shuffle string `json:"shuffle"`
}

// QueueItem holds a queue entry.
type QueueItem struct {
	QID    string `json:"qid"`
	Song   string `json:"song"`
	Artist string `json:"artist"`
	Album  string `json:"album"`
}

// QuickSelect holds a quick select preset.
type QuickSelect struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func playerFromClient(p client.Player) PlayerInfo {
	return PlayerInfo{
		Name:    p.Name,
		PID:     p.PID.String(),
		Model:   p.Model,
		Version: p.Version,
		Network: p.Network,
		Serial:  p.Serial,
	}
}

// ListPlayers returns all available players.
func ListPlayers(c *client.Client) ([]PlayerInfo, error) {
	players, err := c.GetPlayers()
	if err != nil {
		return nil, err
	}
	result := make([]PlayerInfo, len(players))
	for i, p := range players {
		result[i] = playerFromClient(p)
	}
	return result, nil
}

// GetPlayerInfo returns info for a specific player.
func GetPlayerInfo(c *client.Client, pid string) (*PlayerInfo, error) {
	p, err := c.GetPlayerInfo(pid)
	if err != nil {
		return nil, err
	}
	info := playerFromClient(*p)
	return &info, nil
}

// GetPlayState returns the current play state (play/pause/stop).
func GetPlayState(c *client.Client, pid string) (string, error) {
	return c.GetPlayState(pid)
}

// SetPlayState sets the play state (play/pause/stop).
func SetPlayState(c *client.Client, pid, state string) error {
	return c.SetPlayState(pid, state)
}

// GetVolume returns the player volume level.
func GetVolume(c *client.Client, pid string) (int, error) {
	level, err := c.GetVolume(pid)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(level)
}

// SetVolume sets the player volume (0-100).
func SetVolume(c *client.Client, pid string, level int) error {
	return c.SetVolume(pid, fmt.Sprintf("%d", level))
}

// VolumeUp increases the player volume by step.
func VolumeUp(c *client.Client, pid string, step int) error {
	return c.VolumeUp(pid, fmt.Sprintf("%d", step))
}

// VolumeDown decreases the player volume by step.
func VolumeDown(c *client.Client, pid string, step int) error {
	return c.VolumeDown(pid, fmt.Sprintf("%d", step))
}

// GetMute returns the mute state (on/off).
func GetMute(c *client.Client, pid string) (string, error) {
	return c.GetMute(pid)
}

// SetMute sets the mute state (on/off).
func SetMute(c *client.Client, pid, state string) error {
	return c.SetMute(pid, state)
}

// ToggleMute toggles the mute state.
func ToggleMute(c *client.Client, pid string) error {
	return c.ToggleMute(pid)
}

// GetNowPlaying returns currently playing media info.
func GetNowPlaying(c *client.Client, pid string) (*NowPlaying, error) {
	np, err := c.GetNowPlayingMedia(pid)
	if err != nil {
		return nil, err
	}
	return &NowPlaying{
		Type:    np.Type,
		Song:    np.Song,
		Artist:  np.Artist,
		Album:   np.Album,
		Station: np.Station,
	}, nil
}

// GetPlayMode returns the repeat and shuffle settings.
func GetPlayMode(c *client.Client, pid string) (*PlayMode, error) {
	repeat, shuffle, err := c.GetPlayMode(pid)
	if err != nil {
		return nil, err
	}
	return &PlayMode{Repeat: repeat, Shuffle: shuffle}, nil
}

// SetPlayMode sets repeat and shuffle mode.
func SetPlayMode(c *client.Client, pid, repeat, shuffle string) error {
	return c.SetPlayMode(pid, repeat, shuffle)
}

// GetQueue returns the player queue.
func GetQueue(c *client.Client, pid, rangeStart, rangeEnd string) ([]QueueItem, error) {
	items, err := c.GetQueue(pid, rangeStart, rangeEnd)
	if err != nil {
		return nil, err
	}
	result := make([]QueueItem, len(items))
	for i, item := range items {
		result[i] = QueueItem{
			QID:    item.QID.String(),
			Song:   item.Song,
			Artist: item.Artist,
			Album:  item.Album,
		}
	}
	return result, nil
}

// PlayQueueItem plays a specific queue item.
func PlayQueueItem(c *client.Client, pid, qid string) error {
	return c.PlayQueueItem(pid, qid)
}

// RemoveFromQueue removes item(s) from the queue.
func RemoveFromQueue(c *client.Client, pid, qid string) error {
	return c.RemoveFromQueue(pid, qid)
}

// SaveQueue saves the current queue as a playlist.
func SaveQueue(c *client.Client, pid, name string) error {
	return c.SaveQueue(pid, name)
}

// ClearQueue clears the player queue.
func ClearQueue(c *client.Client, pid string) error {
	return c.ClearQueue(pid)
}

// MoveQueue moves queue items.
func MoveQueue(c *client.Client, pid, sqid, dqid string) error {
	return c.MoveQueue(pid, sqid, dqid)
}

// PlayNext skips to the next track.
func PlayNext(c *client.Client, pid string) error {
	return c.PlayNext(pid)
}

// PlayPrevious goes back to the previous track.
func PlayPrevious(c *client.Client, pid string) error {
	return c.PlayPrevious(pid)
}

// SetQuickSelect sets a quick select preset.
func SetQuickSelect(c *client.Client, pid, id string) error {
	return c.SetQuickSelect(pid, id)
}

// PlayQuickSelect plays a quick select preset.
func PlayQuickSelect(c *client.Client, pid, id string) error {
	return c.PlayQuickSelect(pid, id)
}

// GetQuickSelects returns quick select presets.
func GetQuickSelects(c *client.Client, pid, id string) ([]QuickSelect, error) {
	qs, err := c.GetQuickSelects(pid, id)
	if err != nil {
		return nil, err
	}
	result := make([]QuickSelect, len(qs))
	for i, q := range qs {
		result[i] = QuickSelect{ID: q.ID, Name: q.Name}
	}
	return result, nil
}

// CheckFirmwareUpdate checks for firmware updates.
func CheckFirmwareUpdate(c *client.Client, pid string) (string, error) {
	fu, err := c.CheckFirmwareUpdate(pid)
	if err != nil {
		return "", err
	}
	return fu.Update, nil
}
