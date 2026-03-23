package ops

import (
	"encoding/json"

	"github.com/jrogala/heos-cli/client"
)

// MusicSource holds a music source entry.
type MusicSource struct {
	Name      string `json:"name"`
	SID       string `json:"sid"`
	Type      string `json:"type"`
	Available string `json:"available,omitempty"`
}

// BrowseItem holds a browse result item.
type BrowseItem struct {
	Name   string `json:"name"`
	Type   string `json:"type,omitempty"`
	CID    string `json:"cid,omitempty"`
	MID    string `json:"mid,omitempty"`
	Artist string `json:"artist,omitempty"`
	Album  string `json:"album,omitempty"`
}

// SearchCriteria holds a search criteria entry.
type SearchCriteria struct {
	Name string `json:"name"`
	SCID string `json:"scid"`
}

// GetMusicSources returns available music sources.
func GetMusicSources(c *client.Client) ([]MusicSource, error) {
	sources, err := c.GetMusicSources()
	if err != nil {
		return nil, err
	}
	result := make([]MusicSource, len(sources))
	for i, s := range sources {
		result[i] = MusicSource{
			Name:      s.Name,
			SID:       s.SID.String(),
			Type:      s.Type,
			Available: s.Available,
		}
	}
	return result, nil
}

// GetSourceInfo returns info for a specific source.
func GetSourceInfo(c *client.Client, sid string) ([]MusicSource, error) {
	sources, err := c.GetSourceInfo(sid)
	if err != nil {
		return nil, err
	}
	result := make([]MusicSource, len(sources))
	for i, s := range sources {
		result[i] = MusicSource{
			Name:      s.Name,
			SID:       s.SID.String(),
			Type:      s.Type,
			Available: s.Available,
		}
	}
	return result, nil
}

// BrowseSource browses a music source.
func BrowseSource(c *client.Client, sid, rangeStart, rangeEnd string) ([]BrowseItem, error) {
	resp, err := c.BrowseSource(sid, rangeStart, rangeEnd)
	if err != nil {
		return nil, err
	}
	return parseBrowseItems(resp.Payload)
}

// BrowseSourceContainers browses a container within a source.
func BrowseSourceContainers(c *client.Client, sid, cid, rangeStart, rangeEnd string) ([]BrowseItem, error) {
	resp, err := c.BrowseSourceContainers(sid, cid, rangeStart, rangeEnd)
	if err != nil {
		return nil, err
	}
	return parseBrowseItems(resp.Payload)
}

// GetSearchCriteria returns search criteria for a source.
func GetSearchCriteria(c *client.Client, sid string) ([]SearchCriteria, error) {
	criteria, err := c.GetSearchCriteria(sid)
	if err != nil {
		return nil, err
	}
	result := make([]SearchCriteria, len(criteria))
	for i, sc := range criteria {
		result[i] = SearchCriteria{Name: sc.Name, SCID: sc.SCID}
	}
	return result, nil
}

// Search searches a music source.
func Search(c *client.Client, sid, search, scid, rangeStart, rangeEnd string) ([]BrowseItem, error) {
	resp, err := c.Search(sid, search, scid, rangeStart, rangeEnd)
	if err != nil {
		return nil, err
	}
	return parseBrowseItems(resp.Payload)
}

// PlayStation plays a station.
func PlayStation(c *client.Client, pid, sid, cid, mid, name string) error {
	return c.PlayStation(pid, sid, cid, mid, name)
}

// PlayPreset plays a preset.
func PlayPreset(c *client.Client, pid, preset string) error {
	return c.PlayPreset(pid, preset)
}

// PlayInput switches to a physical input.
func PlayInput(c *client.Client, pid, input, spid string) error {
	return c.PlayInput(pid, input, spid)
}

// PlayURL plays a URL on a player.
func PlayURL(c *client.Client, pid, url string) error {
	return c.PlayURL(pid, url)
}

// AddToQueue adds a container to the queue.
func AddToQueue(c *client.Client, pid, sid, cid, aid string) error {
	return c.AddToQueue(pid, sid, cid, aid)
}

// AddTrackToQueue adds a track to the queue.
func AddTrackToQueue(c *client.Client, pid, sid, cid, mid, aid string) error {
	return c.AddTrackToQueue(pid, sid, cid, mid, aid)
}

// RenamePlaylist renames a playlist.
func RenamePlaylist(c *client.Client, sid, cid, name string) error {
	return c.RenamePlaylist(sid, cid, name)
}

// DeletePlaylist deletes a playlist.
func DeletePlaylist(c *client.Client, sid, cid string) error {
	return c.DeletePlaylist(sid, cid)
}

// SetServiceOption sets a service option.
func SetServiceOption(c *client.Client, args map[string]string) error {
	return c.SetServiceOption(args)
}

func parseBrowseItems(payload json.RawMessage) ([]BrowseItem, error) {
	if payload == nil {
		return nil, nil
	}
	var items []client.BrowseItem
	if err := json.Unmarshal(payload, &items); err != nil {
		return nil, err
	}
	result := make([]BrowseItem, len(items))
	for i, item := range items {
		result[i] = BrowseItem{
			Name:   item.Name,
			Type:   item.Type,
			CID:    item.CID,
			MID:    item.MID,
			Artist: item.Artist,
			Album:  item.Album,
		}
	}
	return result, nil
}
