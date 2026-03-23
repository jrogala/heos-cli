package ops

import (
	"fmt"
	"strconv"

	"github.com/jrogala/heos-cli/client"
)

// GroupInfo holds group details.
type GroupInfo struct {
	Name    string              `json:"name"`
	GID     string              `json:"gid"`
	Players []GroupPlayerInfo   `json:"players"`
}

// GroupPlayerInfo holds a player within a group.
type GroupPlayerInfo struct {
	Name string `json:"name"`
	PID  string `json:"pid"`
	Role string `json:"role"`
}

func groupFromClient(g client.Group) GroupInfo {
	players := make([]GroupPlayerInfo, len(g.Players))
	for i, p := range g.Players {
		players[i] = GroupPlayerInfo{
			Name: p.Name,
			PID:  p.PID.String(),
			Role: p.Role,
		}
	}
	return GroupInfo{
		Name:    g.Name,
		GID:     g.GID.String(),
		Players: players,
	}
}

// ListGroups returns all speaker groups.
func ListGroups(c *client.Client) ([]GroupInfo, error) {
	groups, err := c.GetGroups()
	if err != nil {
		return nil, err
	}
	result := make([]GroupInfo, len(groups))
	for i, g := range groups {
		result[i] = groupFromClient(g)
	}
	return result, nil
}

// GetGroupInfo returns info for a specific group.
func GetGroupInfo(c *client.Client, gid string) (*GroupInfo, error) {
	g, err := c.GetGroupInfo(gid)
	if err != nil {
		return nil, err
	}
	info := groupFromClient(*g)
	return &info, nil
}

// SetGroup creates, modifies, or ungroups players.
type SetGroupResult struct {
	GID string `json:"gid,omitempty"`
}

// SetGroup sets a group (first PID is leader, single PID ungroups).
func SetGroup(c *client.Client, pid string) (*SetGroupResult, error) {
	msg, err := c.SetGroup(pid)
	if err != nil {
		return nil, err
	}
	return &SetGroupResult{GID: msg["gid"]}, nil
}

// GetGroupVolume returns the group volume level.
func GetGroupVolume(c *client.Client, gid string) (int, error) {
	level, err := c.GetGroupVolume(gid)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(level)
}

// SetGroupVolume sets the group volume (0-100).
func SetGroupVolume(c *client.Client, gid string, level int) error {
	return c.SetGroupVolume(gid, fmt.Sprintf("%d", level))
}

// GroupVolumeUp increases group volume by step.
func GroupVolumeUp(c *client.Client, gid string, step int) error {
	return c.GroupVolumeUp(gid, fmt.Sprintf("%d", step))
}

// GroupVolumeDown decreases group volume by step.
func GroupVolumeDown(c *client.Client, gid string, step int) error {
	return c.GroupVolumeDown(gid, fmt.Sprintf("%d", step))
}

// GetGroupMute returns the group mute state (on/off).
func GetGroupMute(c *client.Client, gid string) (string, error) {
	return c.GetGroupMute(gid)
}

// SetGroupMute sets the group mute state (on/off).
func SetGroupMute(c *client.Client, gid, state string) error {
	return c.SetGroupMute(gid, state)
}

// ToggleGroupMute toggles group mute.
func ToggleGroupMute(c *client.Client, gid string) error {
	return c.ToggleGroupMute(gid)
}
