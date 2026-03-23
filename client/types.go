package client

import "encoding/json"

// Response is the top-level JSON envelope from HEOS.
type Response struct {
	HEOS    HeosHeader       `json:"heos"`
	Payload json.RawMessage  `json:"payload,omitempty"`
	Options json.RawMessage  `json:"options,omitempty"`
}

// HeosHeader is the "heos" object in every response.
type HeosHeader struct {
	Command string `json:"command"`
	Result  string `json:"result"`
	Message string `json:"message"`
}

// Player represents a HEOS player.
type Player struct {
	Name     string          `json:"name"`
	PID      json.Number     `json:"pid"`
	GID      json.Number     `json:"gid,omitempty"`
	Model    string          `json:"model"`
	Version  string          `json:"version"`
	Network  string          `json:"network"`
	Lineout  json.Number     `json:"lineout,omitempty"`
	Control  json.Number     `json:"control,omitempty"`
	Serial   string          `json:"serial,omitempty"`
}

// Group represents a HEOS speaker group.
type Group struct {
	Name    string        `json:"name"`
	GID     json.Number   `json:"gid"`
	Players []GroupPlayer `json:"players"`
}

// GroupPlayer is a player within a group.
type GroupPlayer struct {
	Name string      `json:"name"`
	PID  json.Number `json:"pid"`
	Role string      `json:"role"`
}

// FlexString handles JSON fields that can be either a string or a number.
type FlexString string

func (f *FlexString) UnmarshalJSON(data []byte) error {
	// Try string first
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*f = FlexString(s)
		return nil
	}
	// Try number
	var n json.Number
	if err := json.Unmarshal(data, &n); err == nil {
		*f = FlexString(n.String())
		return nil
	}
	*f = FlexString(string(data))
	return nil
}

func (f FlexString) String() string { return string(f) }

// NowPlaying represents currently playing media.
type NowPlaying struct {
	Type     string      `json:"type"`
	Song     string      `json:"song"`
	Album    string      `json:"album"`
	Artist   string      `json:"artist"`
	ImageURL string      `json:"image_url"`
	MID      FlexString  `json:"mid"`
	QID      FlexString  `json:"qid"`
	SID      json.Number `json:"sid"`
	AlbumID  string      `json:"album_id,omitempty"`
	Station  string      `json:"station,omitempty"`
}

// QueueItem represents an item in the player queue.
type QueueItem struct {
	Song     string     `json:"song"`
	Album    string     `json:"album"`
	Artist   string     `json:"artist"`
	ImageURL string     `json:"image_url"`
	QID      FlexString `json:"qid"`
	MID      FlexString `json:"mid"`
	AlbumID  string     `json:"album_id,omitempty"`
}

// MusicSource represents a music source from browse/get_music_sources.
type MusicSource struct {
	Name            string      `json:"name"`
	ImageURL        string      `json:"image_url"`
	Type            string      `json:"type"`
	SID             json.Number `json:"sid"`
	Available       string      `json:"available,omitempty"`
	ServiceUsername string      `json:"service_username,omitempty"`
}

// BrowseItem represents an item returned by browse commands.
type BrowseItem struct {
	Container string      `json:"container,omitempty"`
	Playable  string      `json:"playable,omitempty"`
	Type      string      `json:"type,omitempty"`
	Name      string      `json:"name"`
	ImageURL  string      `json:"image_url,omitempty"`
	Artist    string      `json:"artist,omitempty"`
	Album     string      `json:"album,omitempty"`
	CID       string      `json:"cid,omitempty"`
	MID       string      `json:"mid,omitempty"`
	SID       json.Number `json:"sid,omitempty"`
}

// SearchCriteria represents a search criteria from get_search_criteria.
type SearchCriteria struct {
	Name     string `json:"name"`
	SCID     string `json:"scid"`
	Wildcard string `json:"wildcard,omitempty"`
	Playable string `json:"playable,omitempty"`
	CID      string `json:"cid,omitempty"`
}

// QuickSelect represents a quick select preset.
type QuickSelect struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// FirmwareUpdate represents firmware update status.
type FirmwareUpdate struct {
	Update string `json:"update"`
}

// AlbumMetadata represents album metadata from retrieve_metadata.
type AlbumMetadata struct {
	AlbumID string          `json:"album_id"`
	Images  []AlbumImage    `json:"images"`
}

// AlbumImage represents an image in album metadata.
type AlbumImage struct {
	ImageURL string      `json:"image_url"`
	Width    json.Number `json:"width,omitempty"`
}
