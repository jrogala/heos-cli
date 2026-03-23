package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"
)

// Client communicates with a HEOS speaker over TCP.
type Client struct {
	host   string
	port   int
	conn   net.Conn
	reader *bufio.Reader
}

// New creates a new Client (does not connect yet).
func New(host string, port int) *Client {
	return &Client{host: host, port: port}
}

// wake sends an SSDP M-SEARCH to wake the HEOS CLI module from dormant mode.
func wake() {
	msg := "M-SEARCH * HTTP/1.1\r\n" +
		"HOST: 239.255.255.250:1900\r\n" +
		"MAN: \"ssdp:discover\"\r\n" +
		"MX: 3\r\n" +
		"ST: urn:schemas-denon-com:device:ACT-Denon:1\r\n\r\n"
	conn, err := net.DialUDP("udp4", nil, &net.UDPAddr{IP: net.IPv4(239, 255, 255, 250), Port: 1900})
	if err != nil {
		return
	}
	defer conn.Close()
	conn.Write([]byte(msg))
}

// Connect establishes a TCP connection to the HEOS speaker.
// Retries with SSDP wake since the CLI module may be in dormant mode.
func (c *Client) Connect() error {
	addr := fmt.Sprintf("%s:%d", c.host, c.port)
	var lastErr error
	for attempt := 0; attempt < 5; attempt++ {
		if attempt > 0 {
			wake()
			time.Sleep(2 * time.Second)
		}
		conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
		if err == nil {
			c.conn = conn
			c.reader = bufio.NewReader(conn)
			return nil
		}
		lastErr = err
	}
	return fmt.Errorf("connecting to %s (speaker may be in dormant mode): %w", addr, lastErr)
}

// Close closes the TCP connection.
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// buildCommand constructs a HEOS command string.
func buildCommand(group, cmd string, args map[string]string) string {
	s := fmt.Sprintf("heos://%s/%s", group, cmd)
	if len(args) > 0 {
		pairs := make([]string, 0, len(args))
		for k, v := range args {
			pairs = append(pairs, k+"="+v)
		}
		s += "?" + strings.Join(pairs, "&")
	}
	return s + "\r\n"
}

// do sends a command and reads the JSON response, skipping event messages.
func (c *Client) do(group, cmd string, args map[string]string) (*Response, error) {
	command := buildCommand(group, cmd, args)
	if _, err := c.conn.Write([]byte(command)); err != nil {
		return nil, fmt.Errorf("sending command: %w", err)
	}

	for {
		c.conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		line, err := c.reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("reading response: %w", err)
		}
		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			continue
		}

		var resp Response
		if err := json.Unmarshal([]byte(line), &resp); err != nil {
			return nil, fmt.Errorf("parsing response: %w", err)
		}

		// Skip unsolicited event messages
		if strings.HasPrefix(resp.HEOS.Command, "event/") {
			continue
		}

		// Check for "command under process" and wait for real response
		if resp.HEOS.Result == "success" && resp.HEOS.Message == "command under process" {
			continue
		}

		// Check for errors
		if resp.HEOS.Result == "fail" {
			msg := ParseMessage(resp.HEOS.Message)
			return nil, parseHEOSError(msg)
		}

		return &resp, nil
	}
}

// --- System Commands ---

func (c *Client) RegisterForChangeEvents(enable bool) error {
	val := "off"
	if enable {
		val = "on"
	}
	_, err := c.do("system", "register_for_change_events", map[string]string{"enable": val})
	return err
}

func (c *Client) CheckAccount() (map[string]string, error) {
	resp, err := c.do("system", "check_account", nil)
	if err != nil {
		return nil, err
	}
	return ParseMessage(resp.HEOS.Message), nil
}

func (c *Client) SignIn(username, password string) error {
	_, err := c.do("system", "sign_in", map[string]string{
		"un": EncodeHEOS(username),
		"pw": EncodeHEOS(password),
	})
	return err
}

func (c *Client) SignOut() error {
	_, err := c.do("system", "sign_out", nil)
	return err
}

func (c *Client) HeartBeat() error {
	_, err := c.do("system", "heart_beat", nil)
	return err
}

func (c *Client) Reboot() error {
	_, err := c.do("system", "reboot", nil)
	return err
}

func (c *Client) SetPrettyJSON(enable bool) error {
	val := "off"
	if enable {
		val = "on"
	}
	_, err := c.do("system", "prettify_json_response", map[string]string{"enable": val})
	return err
}

// --- Player Commands ---

func (c *Client) GetPlayers() ([]Player, error) {
	resp, err := c.do("player", "get_players", nil)
	if err != nil {
		return nil, err
	}
	var players []Player
	if resp.Payload != nil {
		if err := json.Unmarshal(resp.Payload, &players); err != nil {
			return nil, fmt.Errorf("parsing players: %w", err)
		}
	}
	return players, nil
}

func (c *Client) GetPlayerInfo(pid string) (*Player, error) {
	resp, err := c.do("player", "get_player_info", map[string]string{"pid": pid})
	if err != nil {
		return nil, err
	}
	var player Player
	if resp.Payload != nil {
		if err := json.Unmarshal(resp.Payload, &player); err != nil {
			return nil, fmt.Errorf("parsing player info: %w", err)
		}
	}
	return &player, nil
}

func (c *Client) GetPlayState(pid string) (string, error) {
	resp, err := c.do("player", "get_play_state", map[string]string{"pid": pid})
	if err != nil {
		return "", err
	}
	msg := ParseMessage(resp.HEOS.Message)
	return msg["state"], nil
}

func (c *Client) SetPlayState(pid, state string) error {
	_, err := c.do("player", "set_play_state", map[string]string{"pid": pid, "state": state})
	return err
}

func (c *Client) GetNowPlayingMedia(pid string) (*NowPlaying, error) {
	resp, err := c.do("player", "get_now_playing_media", map[string]string{"pid": pid})
	if err != nil {
		return nil, err
	}
	var np NowPlaying
	if resp.Payload != nil {
		if err := json.Unmarshal(resp.Payload, &np); err != nil {
			return nil, fmt.Errorf("parsing now playing: %w", err)
		}
	}
	return &np, nil
}

func (c *Client) GetVolume(pid string) (string, error) {
	resp, err := c.do("player", "get_volume", map[string]string{"pid": pid})
	if err != nil {
		return "", err
	}
	msg := ParseMessage(resp.HEOS.Message)
	return msg["level"], nil
}

func (c *Client) SetVolume(pid, level string) error {
	_, err := c.do("player", "set_volume", map[string]string{"pid": pid, "level": level})
	return err
}

func (c *Client) VolumeUp(pid, step string) error {
	_, err := c.do("player", "volume_up", map[string]string{"pid": pid, "step": step})
	return err
}

func (c *Client) VolumeDown(pid, step string) error {
	_, err := c.do("player", "volume_down", map[string]string{"pid": pid, "step": step})
	return err
}

func (c *Client) GetMute(pid string) (string, error) {
	resp, err := c.do("player", "get_mute", map[string]string{"pid": pid})
	if err != nil {
		return "", err
	}
	msg := ParseMessage(resp.HEOS.Message)
	return msg["state"], nil
}

func (c *Client) SetMute(pid, state string) error {
	_, err := c.do("player", "set_mute", map[string]string{"pid": pid, "state": state})
	return err
}

func (c *Client) ToggleMute(pid string) error {
	_, err := c.do("player", "toggle_mute", map[string]string{"pid": pid})
	return err
}

func (c *Client) GetPlayMode(pid string) (repeat string, shuffle string, err error) {
	resp, err := c.do("player", "get_play_mode", map[string]string{"pid": pid})
	if err != nil {
		return "", "", err
	}
	msg := ParseMessage(resp.HEOS.Message)
	return msg["repeat"], msg["shuffle"], nil
}

func (c *Client) SetPlayMode(pid, repeat, shuffle string) error {
	_, err := c.do("player", "set_play_mode", map[string]string{
		"pid": pid, "repeat": repeat, "shuffle": shuffle,
	})
	return err
}

func (c *Client) GetQueue(pid string, rangeStart, rangeEnd string) ([]QueueItem, error) {
	args := map[string]string{"pid": pid}
	if rangeStart != "" && rangeEnd != "" {
		args["range"] = rangeStart + "," + rangeEnd
	}
	resp, err := c.do("player", "get_queue", args)
	if err != nil {
		return nil, err
	}
	var items []QueueItem
	if resp.Payload != nil {
		if err := json.Unmarshal(resp.Payload, &items); err != nil {
			return nil, fmt.Errorf("parsing queue: %w", err)
		}
	}
	return items, nil
}

func (c *Client) PlayQueueItem(pid, qid string) error {
	_, err := c.do("player", "play_queue", map[string]string{"pid": pid, "qid": qid})
	return err
}

func (c *Client) RemoveFromQueue(pid, qid string) error {
	_, err := c.do("player", "remove_from_queue", map[string]string{"pid": pid, "qid": qid})
	return err
}

func (c *Client) SaveQueue(pid, name string) error {
	_, err := c.do("player", "save_queue", map[string]string{"pid": pid, "name": EncodeHEOS(name)})
	return err
}

func (c *Client) ClearQueue(pid string) error {
	_, err := c.do("player", "clear_queue", map[string]string{"pid": pid})
	return err
}

func (c *Client) MoveQueue(pid, sqid, dqid string) error {
	_, err := c.do("player", "move_queue_item", map[string]string{
		"pid": pid, "sqid": sqid, "dqid": dqid,
	})
	return err
}

func (c *Client) PlayNext(pid string) error {
	_, err := c.do("player", "play_next", map[string]string{"pid": pid})
	return err
}

func (c *Client) PlayPrevious(pid string) error {
	_, err := c.do("player", "play_previous", map[string]string{"pid": pid})
	return err
}

func (c *Client) SetQuickSelect(pid, id string) error {
	_, err := c.do("player", "set_quickselect", map[string]string{"pid": pid, "id": id})
	return err
}

func (c *Client) PlayQuickSelect(pid, id string) error {
	_, err := c.do("player", "play_quickselect", map[string]string{"pid": pid, "id": id})
	return err
}

func (c *Client) GetQuickSelects(pid string, id string) ([]QuickSelect, error) {
	args := map[string]string{"pid": pid}
	if id != "" {
		args["id"] = id
	}
	resp, err := c.do("player", "get_quickselects", args)
	if err != nil {
		return nil, err
	}
	var qs []QuickSelect
	if resp.Payload != nil {
		if err := json.Unmarshal(resp.Payload, &qs); err != nil {
			return nil, fmt.Errorf("parsing quickselects: %w", err)
		}
	}
	return qs, nil
}

func (c *Client) CheckFirmwareUpdate(pid string) (*FirmwareUpdate, error) {
	resp, err := c.do("player", "check_update", map[string]string{"pid": pid})
	if err != nil {
		return nil, err
	}
	var fu FirmwareUpdate
	if resp.Payload != nil {
		if err := json.Unmarshal(resp.Payload, &fu); err != nil {
			return nil, fmt.Errorf("parsing firmware update: %w", err)
		}
	}
	return &fu, nil
}

// --- Group Commands ---

func (c *Client) GetGroups() ([]Group, error) {
	resp, err := c.do("group", "get_groups", nil)
	if err != nil {
		return nil, err
	}
	var groups []Group
	if resp.Payload != nil {
		if err := json.Unmarshal(resp.Payload, &groups); err != nil {
			return nil, fmt.Errorf("parsing groups: %w", err)
		}
	}
	return groups, nil
}

func (c *Client) GetGroupInfo(gid string) (*Group, error) {
	resp, err := c.do("group", "get_group_info", map[string]string{"gid": gid})
	if err != nil {
		return nil, err
	}
	var group Group
	if resp.Payload != nil {
		if err := json.Unmarshal(resp.Payload, &group); err != nil {
			return nil, fmt.Errorf("parsing group info: %w", err)
		}
	}
	return &group, nil
}

func (c *Client) SetGroup(pid string) (map[string]string, error) {
	resp, err := c.do("group", "set_group", map[string]string{"pid": pid})
	if err != nil {
		return nil, err
	}
	return ParseMessage(resp.HEOS.Message), nil
}

func (c *Client) GetGroupVolume(gid string) (string, error) {
	resp, err := c.do("group", "get_volume", map[string]string{"gid": gid})
	if err != nil {
		return "", err
	}
	msg := ParseMessage(resp.HEOS.Message)
	return msg["level"], nil
}

func (c *Client) SetGroupVolume(gid, level string) error {
	_, err := c.do("group", "set_volume", map[string]string{"gid": gid, "level": level})
	return err
}

func (c *Client) GroupVolumeUp(gid, step string) error {
	_, err := c.do("group", "volume_up", map[string]string{"gid": gid, "step": step})
	return err
}

func (c *Client) GroupVolumeDown(gid, step string) error {
	_, err := c.do("group", "volume_down", map[string]string{"gid": gid, "step": step})
	return err
}

func (c *Client) GetGroupMute(gid string) (string, error) {
	resp, err := c.do("group", "get_mute", map[string]string{"gid": gid})
	if err != nil {
		return "", err
	}
	msg := ParseMessage(resp.HEOS.Message)
	return msg["state"], nil
}

func (c *Client) SetGroupMute(gid, state string) error {
	_, err := c.do("group", "set_mute", map[string]string{"gid": gid, "state": state})
	return err
}

func (c *Client) ToggleGroupMute(gid string) error {
	_, err := c.do("group", "toggle_mute", map[string]string{"gid": gid})
	return err
}

// --- Browse Commands ---

func (c *Client) GetMusicSources() ([]MusicSource, error) {
	resp, err := c.do("browse", "get_music_sources", nil)
	if err != nil {
		return nil, err
	}
	var sources []MusicSource
	if resp.Payload != nil {
		if err := json.Unmarshal(resp.Payload, &sources); err != nil {
			return nil, fmt.Errorf("parsing music sources: %w", err)
		}
	}
	return sources, nil
}

func (c *Client) GetSourceInfo(sid string) ([]MusicSource, error) {
	resp, err := c.do("browse", "get_source_info", map[string]string{"sid": sid})
	if err != nil {
		return nil, err
	}
	var sources []MusicSource
	if resp.Payload != nil {
		if err := json.Unmarshal(resp.Payload, &sources); err != nil {
			return nil, fmt.Errorf("parsing source info: %w", err)
		}
	}
	return sources, nil
}

func (c *Client) BrowseSource(sid string, rangeStart, rangeEnd string) (*Response, error) {
	args := map[string]string{"sid": sid}
	if rangeStart != "" && rangeEnd != "" {
		args["range"] = rangeStart + "," + rangeEnd
	}
	return c.do("browse", "browse", args)
}

func (c *Client) BrowseSourceContainers(sid, cid string, rangeStart, rangeEnd string) (*Response, error) {
	args := map[string]string{"sid": sid, "cid": cid}
	if rangeStart != "" && rangeEnd != "" {
		args["range"] = rangeStart + "," + rangeEnd
	}
	return c.do("browse", "browse", args)
}

func (c *Client) GetSearchCriteria(sid string) ([]SearchCriteria, error) {
	resp, err := c.do("browse", "get_search_criteria", map[string]string{"sid": sid})
	if err != nil {
		return nil, err
	}
	var criteria []SearchCriteria
	if resp.Payload != nil {
		if err := json.Unmarshal(resp.Payload, &criteria); err != nil {
			return nil, fmt.Errorf("parsing search criteria: %w", err)
		}
	}
	return criteria, nil
}

func (c *Client) Search(sid, search, scid string, rangeStart, rangeEnd string) (*Response, error) {
	args := map[string]string{
		"sid":    sid,
		"search": EncodeHEOS(search),
		"scid":  scid,
	}
	if rangeStart != "" && rangeEnd != "" {
		args["range"] = rangeStart + "," + rangeEnd
	}
	return c.do("browse", "search", args)
}

func (c *Client) PlayStation(pid, sid, cid, mid, name string) error {
	_, err := c.do("browse", "play_stream", map[string]string{
		"pid":  pid,
		"sid":  sid,
		"cid":  cid,
		"mid":  mid,
		"name": EncodeHEOS(name),
	})
	return err
}

func (c *Client) PlayPreset(pid, preset string) error {
	_, err := c.do("browse", "play_preset", map[string]string{"pid": pid, "preset": preset})
	return err
}

func (c *Client) PlayInput(pid, input string, spid string) error {
	args := map[string]string{"pid": pid, "input": input}
	if spid != "" {
		args["spid"] = spid
	}
	_, err := c.do("browse", "play_input", args)
	return err
}

func (c *Client) PlayURL(pid, url string) error {
	_, err := c.do("browse", "play_stream", map[string]string{"pid": pid, "url": url})
	return err
}

func (c *Client) AddToQueue(pid, sid, cid, aid string) error {
	_, err := c.do("browse", "add_to_queue", map[string]string{
		"pid": pid, "sid": sid, "cid": cid, "aid": aid,
	})
	return err
}

func (c *Client) AddTrackToQueue(pid, sid, cid, mid, aid string) error {
	_, err := c.do("browse", "add_to_queue", map[string]string{
		"pid": pid, "sid": sid, "cid": cid, "mid": mid, "aid": aid,
	})
	return err
}

func (c *Client) RenamePlaylist(sid, cid, name string) error {
	_, err := c.do("browse", "rename_playlist", map[string]string{
		"sid": sid, "cid": cid, "name": EncodeHEOS(name),
	})
	return err
}

func (c *Client) DeletePlaylist(sid, cid string) error {
	_, err := c.do("browse", "delete_playlist", map[string]string{"sid": sid, "cid": cid})
	return err
}

func (c *Client) RetrieveMetadata(sid, cid string) ([]AlbumMetadata, error) {
	resp, err := c.do("browse", "retrieve_metadata", map[string]string{"sid": sid, "cid": cid})
	if err != nil {
		return nil, err
	}
	var meta []AlbumMetadata
	if resp.Payload != nil {
		if err := json.Unmarshal(resp.Payload, &meta); err != nil {
			return nil, fmt.Errorf("parsing metadata: %w", err)
		}
	}
	return meta, nil
}

func (c *Client) SetServiceOption(args map[string]string) error {
	_, err := c.do("browse", "set_service_option", args)
	return err
}
