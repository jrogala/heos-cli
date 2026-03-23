package tests

import (
	"context"
	"fmt"
	"strings"

	"github.com/cucumber/godog"
)

func initializeScenario(ctx *godog.ScenarioContext) {
	var sc *scenarioCtx

	ctx.Before(func(ctx context.Context, sc2 *godog.Scenario) (context.Context, error) {
		globalServer.Reset()
		sc = newScenarioCtx()
		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc2 *godog.Scenario, err error) (context.Context, error) {
		if sc != nil {
			sc.cleanup()
		}
		return ctx, nil
	})

	// Connection
	ctx.Step(`^I am connected to a HEOS speaker$`, func() error {
		return sc.connect()
	})

	// System steps
	ctx.Step(`^the speaker responds to heartbeat$`, func() error {
		sc.server.OnSuccess("system/heart_beat", "")
		return nil
	})
	ctx.Step(`^I send a heartbeat$`, func() error {
		sc.lastErr = sc.client.HeartBeat()
		return nil
	})
	ctx.Step(`^the speaker has account "([^"]*)" signed in$`, func(un string) error {
		sc.server.OnSuccess("system/check_account", "signed_in&un="+un)
		return nil
	})
	ctx.Step(`^the speaker has no account signed in$`, func() error {
		sc.server.OnSuccess("system/check_account", "signed_out")
		return nil
	})
	ctx.Step(`^I check the account$`, func() error {
		sc.account, sc.lastErr = sc.client.CheckAccount()
		return nil
	})
	ctx.Step(`^I should see username "([^"]*)"$`, func(expected string) error {
		if sc.account["un"] != expected {
			return fmt.Errorf("expected username %q, got %q", expected, sc.account["un"])
		}
		return nil
	})
	ctx.Step(`^I should see signed out$`, func() error {
		if _, ok := sc.account["un"]; ok {
			return fmt.Errorf("expected signed out, but got username %q", sc.account["un"])
		}
		return nil
	})
	ctx.Step(`^the speaker accepts sign-in$`, func() error {
		sc.server.OnSuccess("system/sign_in", "signed_in&un=test@example.com")
		return nil
	})
	ctx.Step(`^the speaker rejects sign-in$`, func() error {
		sc.server.OnFail("system/sign_in", 6, "Invalid_Credentials")
		return nil
	})
	ctx.Step(`^I sign in with "([^"]*)" and "([^"]*)"$`, func(un, pw string) error {
		sc.lastErr = sc.client.SignIn(un, pw)
		return nil
	})
	ctx.Step(`^I sign out$`, func() error {
		sc.server.OnSuccess("system/sign_out", "signed_out")
		sc.lastErr = sc.client.SignOut()
		return nil
	})

	// Player steps
	ctx.Step(`^the speaker has (\d+) players?$`, func(n int) error {
		players := make([]map[string]any, n)
		for i := 0; i < n; i++ {
			players[i] = map[string]any{
				"name":    fmt.Sprintf("Player %d", i+1),
				"pid":     i + 1,
				"model":   "Denon Home 150",
				"version": "3.34.510",
				"network": "wifi",
				"serial":  fmt.Sprintf("SERIAL%d", i+1),
			}
		}
		sc.server.OnSuccessWithPayload("player/get_players", "", players)
		return nil
	})
	ctx.Step(`^I list players$`, func() error {
		sc.players, sc.lastErr = sc.client.GetPlayers()
		return nil
	})
	ctx.Step(`^I should get (\d+) players?$`, func(expected int) error {
		if len(sc.players) != expected {
			return fmt.Errorf("expected %d players, got %d", expected, len(sc.players))
		}
		return nil
	})
	ctx.Step(`^the speaker has player "([^"]*)" with pid (\d+)$`, func(name string, pid int) error {
		player := map[string]any{
			"name":    name,
			"pid":     pid,
			"model":   "Denon Home 150",
			"version": "3.34.510",
			"network": "wifi",
			"serial":  "AABBCC",
		}
		sc.server.OnSuccessWithPayload("player/get_player_info", fmt.Sprintf("pid=%d", pid), player)
		return nil
	})
	ctx.Step(`^I get info for player (\d+)$`, func(pid int) error {
		sc.player, sc.lastErr = sc.client.GetPlayerInfo(fmt.Sprintf("%d", pid))
		return nil
	})
	ctx.Step(`^the player name should be "([^"]*)"$`, func(expected string) error {
		if sc.player == nil {
			return fmt.Errorf("no player info")
		}
		if sc.player.Name != expected {
			return fmt.Errorf("expected name %q, got %q", expected, sc.player.Name)
		}
		return nil
	})

	// Volume steps
	ctx.Step(`^player (\d+) has volume (\d+)$`, func(pid, level int) error {
		sc.server.OnSuccess("player/get_volume", fmt.Sprintf("pid=%d&level=%d", pid, level))
		return nil
	})
	ctx.Step(`^I get volume for player (\d+)$`, func(pid int) error {
		sc.stringVal, sc.lastErr = sc.client.GetVolume(fmt.Sprintf("%d", pid))
		return nil
	})
	ctx.Step(`^the volume should be "([^"]*)"$`, func(expected string) error {
		if sc.stringVal != expected {
			return fmt.Errorf("expected volume %q, got %q", expected, sc.stringVal)
		}
		return nil
	})
	ctx.Step(`^the speaker accepts volume changes$`, func() error {
		sc.server.OnSuccess("player/set_volume", "pid=1&level=50")
		sc.server.OnSuccess("player/volume_up", "pid=1")
		sc.server.OnSuccess("player/volume_down", "pid=1")
		return nil
	})
	ctx.Step(`^I set volume to (\d+) for player (\d+)$`, func(level, pid int) error {
		sc.lastErr = sc.client.SetVolume(fmt.Sprintf("%d", pid), fmt.Sprintf("%d", level))
		return nil
	})
	ctx.Step(`^I increase volume for player (\d+)$`, func(pid int) error {
		sc.lastErr = sc.client.VolumeUp(fmt.Sprintf("%d", pid), "5")
		return nil
	})
	ctx.Step(`^I decrease volume for player (\d+)$`, func(pid int) error {
		sc.lastErr = sc.client.VolumeDown(fmt.Sprintf("%d", pid), "5")
		return nil
	})

	// Play state steps
	ctx.Step(`^player (\d+) is "([^"]*)"$`, func(pid int, state string) error {
		sc.server.OnSuccess("player/get_play_state", fmt.Sprintf("pid=%d&state=%s", pid, state))
		return nil
	})
	ctx.Step(`^I get play state for player (\d+)$`, func(pid int) error {
		sc.stringVal, sc.lastErr = sc.client.GetPlayState(fmt.Sprintf("%d", pid))
		return nil
	})
	ctx.Step(`^the state should be "([^"]*)"$`, func(expected string) error {
		if sc.stringVal != expected {
			return fmt.Errorf("expected state %q, got %q", expected, sc.stringVal)
		}
		return nil
	})
	ctx.Step(`^the speaker accepts play state changes$`, func() error {
		sc.server.OnSuccess("player/set_play_state", "pid=1&state=play")
		return nil
	})
	ctx.Step(`^I set play state "([^"]*)" for player (\d+)$`, func(state string, pid int) error {
		sc.lastErr = sc.client.SetPlayState(fmt.Sprintf("%d", pid), state)
		return nil
	})

	// Mute steps
	ctx.Step(`^player (\d+) is muted$`, func(pid int) error {
		sc.server.OnSuccess("player/get_mute", fmt.Sprintf("pid=%d&state=on", pid))
		return nil
	})
	ctx.Step(`^player (\d+) is not muted$`, func(pid int) error {
		sc.server.OnSuccess("player/get_mute", fmt.Sprintf("pid=%d&state=off", pid))
		return nil
	})
	ctx.Step(`^I get mute state for player (\d+)$`, func(pid int) error {
		sc.stringVal, sc.lastErr = sc.client.GetMute(fmt.Sprintf("%d", pid))
		return nil
	})
	ctx.Step(`^the mute state should be "([^"]*)"$`, func(expected string) error {
		if sc.stringVal != expected {
			return fmt.Errorf("expected mute %q, got %q", expected, sc.stringVal)
		}
		return nil
	})
	ctx.Step(`^the speaker accepts mute changes$`, func() error {
		sc.server.OnSuccess("player/set_mute", "pid=1")
		sc.server.OnSuccess("player/toggle_mute", "pid=1")
		return nil
	})
	ctx.Step(`^I mute player (\d+)$`, func(pid int) error {
		sc.lastErr = sc.client.SetMute(fmt.Sprintf("%d", pid), "on")
		return nil
	})
	ctx.Step(`^I unmute player (\d+)$`, func(pid int) error {
		sc.lastErr = sc.client.SetMute(fmt.Sprintf("%d", pid), "off")
		return nil
	})
	ctx.Step(`^I toggle mute for player (\d+)$`, func(pid int) error {
		sc.lastErr = sc.client.ToggleMute(fmt.Sprintf("%d", pid))
		return nil
	})

	// Now playing steps
	ctx.Step(`^player (\d+) is playing "([^"]*)" by "([^"]*)" from "([^"]*)"$`, func(pid int, song, artist, album string) error {
		np := map[string]any{
			"type":      "song",
			"song":      song,
			"artist":    artist,
			"album":     album,
			"image_url": "http://example.com/art.jpg",
			"mid":       "1",
			"qid":       "1",
			"sid":       1,
		}
		sc.server.OnSuccessWithPayload("player/get_now_playing_media", fmt.Sprintf("pid=%d", pid), np)
		return nil
	})
	ctx.Step(`^I get now playing for player (\d+)$`, func(pid int) error {
		sc.nowPlaying, sc.lastErr = sc.client.GetNowPlayingMedia(fmt.Sprintf("%d", pid))
		return nil
	})
	ctx.Step(`^the song should be "([^"]*)"$`, func(expected string) error {
		if sc.nowPlaying == nil {
			return fmt.Errorf("no now playing info")
		}
		if sc.nowPlaying.Song != expected {
			return fmt.Errorf("expected song %q, got %q", expected, sc.nowPlaying.Song)
		}
		return nil
	})
	ctx.Step(`^the artist should be "([^"]*)"$`, func(expected string) error {
		if sc.nowPlaying.Artist != expected {
			return fmt.Errorf("expected artist %q, got %q", expected, sc.nowPlaying.Artist)
		}
		return nil
	})

	// Play mode steps
	ctx.Step(`^player (\d+) has play mode repeat "([^"]*)" shuffle "([^"]*)"$`, func(pid int, repeat, shuffle string) error {
		sc.server.OnSuccess("player/get_play_mode", fmt.Sprintf("pid=%d&repeat=%s&shuffle=%s", pid, repeat, shuffle))
		return nil
	})
	ctx.Step(`^I get play mode for player (\d+)$`, func(pid int) error {
		sc.repeatVal, sc.shuffleVal, sc.lastErr = sc.client.GetPlayMode(fmt.Sprintf("%d", pid))
		return nil
	})
	ctx.Step(`^the repeat mode should be "([^"]*)"$`, func(expected string) error {
		if sc.repeatVal != expected {
			return fmt.Errorf("expected repeat %q, got %q", expected, sc.repeatVal)
		}
		return nil
	})
	ctx.Step(`^the shuffle mode should be "([^"]*)"$`, func(expected string) error {
		if sc.shuffleVal != expected {
			return fmt.Errorf("expected shuffle %q, got %q", expected, sc.shuffleVal)
		}
		return nil
	})
	ctx.Step(`^the speaker accepts play mode changes$`, func() error {
		sc.server.OnSuccess("player/set_play_mode", "pid=1")
		return nil
	})
	ctx.Step(`^I set play mode repeat "([^"]*)" shuffle "([^"]*)" for player (\d+)$`, func(repeat, shuffle string, pid int) error {
		sc.lastErr = sc.client.SetPlayMode(fmt.Sprintf("%d", pid), repeat, shuffle)
		return nil
	})

	// Queue steps
	ctx.Step(`^player (\d+) has (\d+) items in queue$`, func(pid, n int) error {
		items := make([]map[string]any, n)
		for i := 0; i < n; i++ {
			items[i] = map[string]any{
				"song":      fmt.Sprintf("Song %d", i+1),
				"album":     "Test Album",
				"artist":    "Test Artist",
				"image_url": "http://example.com/art.jpg",
				"qid":       fmt.Sprintf("%d", i+1),
				"mid":       fmt.Sprintf("%d", i+1),
			}
		}
		sc.server.OnSuccessWithPayload("player/get_queue", fmt.Sprintf("pid=%d", pid), items)
		return nil
	})
	ctx.Step(`^I get queue for player (\d+)$`, func(pid int) error {
		sc.queueItems, sc.lastErr = sc.client.GetQueue(fmt.Sprintf("%d", pid), "", "")
		return nil
	})
	ctx.Step(`^I should get (\d+) queue items$`, func(expected int) error {
		if len(sc.queueItems) != expected {
			return fmt.Errorf("expected %d queue items, got %d", expected, len(sc.queueItems))
		}
		return nil
	})
	ctx.Step(`^the speaker accepts queue operations$`, func() error {
		sc.server.OnSuccess("player/play_queue", "pid=1")
		sc.server.OnSuccess("player/remove_from_queue", "pid=1")
		sc.server.OnSuccess("player/clear_queue", "pid=1")
		sc.server.OnSuccess("player/save_queue", "pid=1")
		sc.server.OnSuccess("player/move_queue_item", "pid=1")
		return nil
	})
	ctx.Step(`^I play queue item (\d+) for player (\d+)$`, func(qid, pid int) error {
		sc.lastErr = sc.client.PlayQueueItem(fmt.Sprintf("%d", pid), fmt.Sprintf("%d", qid))
		return nil
	})
	ctx.Step(`^I clear queue for player (\d+)$`, func(pid int) error {
		sc.lastErr = sc.client.ClearQueue(fmt.Sprintf("%d", pid))
		return nil
	})

	// Navigation steps
	ctx.Step(`^the speaker accepts navigation commands$`, func() error {
		sc.server.OnSuccess("player/play_next", "pid=1")
		sc.server.OnSuccess("player/play_previous", "pid=1")
		return nil
	})
	ctx.Step(`^I play next for player (\d+)$`, func(pid int) error {
		sc.lastErr = sc.client.PlayNext(fmt.Sprintf("%d", pid))
		return nil
	})
	ctx.Step(`^I play previous for player (\d+)$`, func(pid int) error {
		sc.lastErr = sc.client.PlayPrevious(fmt.Sprintf("%d", pid))
		return nil
	})

	// Group steps
	ctx.Step(`^the speaker has (\d+) groups?$`, func(n int) error {
		groups := make([]map[string]any, n)
		for i := 0; i < n; i++ {
			groups[i] = map[string]any{
				"name": fmt.Sprintf("Group %d", i+1),
				"gid":  i + 1,
				"players": []map[string]any{
					{"name": "Player 1", "pid": 1, "role": "leader"},
					{"name": "Player 2", "pid": 2, "role": "member"},
				},
			}
		}
		sc.server.OnSuccessWithPayload("group/get_groups", "", groups)
		return nil
	})
	ctx.Step(`^I list groups$`, func() error {
		sc.groups, sc.lastErr = sc.client.GetGroups()
		return nil
	})
	ctx.Step(`^I should get (\d+) groups?$`, func(expected int) error {
		if len(sc.groups) != expected {
			return fmt.Errorf("expected %d groups, got %d", expected, len(sc.groups))
		}
		return nil
	})
	ctx.Step(`^group (\d+) has volume (\d+)$`, func(gid, level int) error {
		sc.server.OnSuccess("group/get_volume", fmt.Sprintf("gid=%d&level=%d", gid, level))
		return nil
	})
	ctx.Step(`^I get volume for group (\d+)$`, func(gid int) error {
		sc.stringVal, sc.lastErr = sc.client.GetGroupVolume(fmt.Sprintf("%d", gid))
		return nil
	})
	ctx.Step(`^group (\d+) is muted$`, func(gid int) error {
		sc.server.OnSuccess("group/get_mute", fmt.Sprintf("gid=%d&state=on", gid))
		return nil
	})
	ctx.Step(`^group (\d+) is not muted$`, func(gid int) error {
		sc.server.OnSuccess("group/get_mute", fmt.Sprintf("gid=%d&state=off", gid))
		return nil
	})
	ctx.Step(`^I get mute state for group (\d+)$`, func(gid int) error {
		sc.stringVal, sc.lastErr = sc.client.GetGroupMute(fmt.Sprintf("%d", gid))
		return nil
	})
	ctx.Step(`^the speaker accepts group volume changes$`, func() error {
		sc.server.OnSuccess("group/set_volume", "gid=1")
		sc.server.OnSuccess("group/volume_up", "gid=1")
		sc.server.OnSuccess("group/volume_down", "gid=1")
		return nil
	})
	ctx.Step(`^I set volume to (\d+) for group (\d+)$`, func(level, gid int) error {
		sc.lastErr = sc.client.SetGroupVolume(fmt.Sprintf("%d", gid), fmt.Sprintf("%d", level))
		return nil
	})
	ctx.Step(`^the speaker accepts group mute changes$`, func() error {
		sc.server.OnSuccess("group/set_mute", "gid=1")
		sc.server.OnSuccess("group/toggle_mute", "gid=1")
		return nil
	})
	ctx.Step(`^I toggle mute for group (\d+)$`, func(gid int) error {
		sc.lastErr = sc.client.ToggleGroupMute(fmt.Sprintf("%d", gid))
		return nil
	})

	// Browse steps
	ctx.Step(`^the speaker has music sources$`, func() error {
		sources := []map[string]any{
			{"name": "Spotify", "image_url": "", "type": "music_service", "sid": 1, "available": "true"},
			{"name": "TuneIn", "image_url": "", "type": "music_service", "sid": 2, "available": "true"},
			{"name": "Local Music", "image_url": "", "type": "heos_server", "sid": 3, "available": "true"},
		}
		sc.server.OnSuccessWithPayload("browse/get_music_sources", "", sources)
		return nil
	})
	ctx.Step(`^I list music sources$`, func() error {
		sc.sources, sc.lastErr = sc.client.GetMusicSources()
		return nil
	})
	ctx.Step(`^I should get (\d+) sources$`, func(expected int) error {
		if len(sc.sources) != expected {
			return fmt.Errorf("expected %d sources, got %d", expected, len(sc.sources))
		}
		return nil
	})
	ctx.Step(`^the speaker accepts play-url$`, func() error {
		sc.server.OnSuccess("browse/play_stream", "pid=1")
		return nil
	})
	ctx.Step(`^I play URL "([^"]*)" on player (\d+)$`, func(url string, pid int) error {
		sc.lastErr = sc.client.PlayURL(fmt.Sprintf("%d", pid), url)
		return nil
	})

	// Error assertions
	ctx.Step(`^the operation should succeed$`, func() error {
		if sc.lastErr != nil {
			return fmt.Errorf("expected success, got error: %v", sc.lastErr)
		}
		return nil
	})
	ctx.Step(`^the operation should fail$`, func() error {
		if sc.lastErr == nil {
			return fmt.Errorf("expected failure, got success")
		}
		return nil
	})
	ctx.Step(`^the error should contain "([^"]*)"$`, func(substr string) error {
		if sc.lastErr == nil {
			return fmt.Errorf("no error to check")
		}
		if !strings.Contains(sc.lastErr.Error(), substr) {
			return fmt.Errorf("error %q does not contain %q", sc.lastErr.Error(), substr)
		}
		return nil
	})
}
