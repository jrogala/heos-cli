// Package ops contains business logic for HEOS operations.
// Functions return Go structs and errors — zero I/O, zero formatting.
package ops

import "github.com/jrogala/heos-cli/client"

// AccountInfo holds HEOS account status.
type AccountInfo struct {
	SignedIn bool   `json:"signed_in"`
	Username string `json:"username,omitempty"`
}

// HeartBeat checks connectivity to the speaker.
func HeartBeat(c *client.Client) error {
	return c.HeartBeat()
}

// CheckAccount returns the current HEOS account status.
func CheckAccount(c *client.Client) (*AccountInfo, error) {
	msg, err := c.CheckAccount()
	if err != nil {
		return nil, err
	}
	un, ok := msg["un"]
	return &AccountInfo{
		SignedIn: ok,
		Username: un,
	}, nil
}

// SignIn authenticates with the HEOS account.
func SignIn(c *client.Client, username, password string) error {
	return c.SignIn(username, password)
}

// SignOut signs out of the HEOS account.
func SignOut(c *client.Client) error {
	return c.SignOut()
}

// Reboot reboots the HEOS speaker.
func Reboot(c *client.Client) error {
	return c.Reboot()
}

// RegisterChangeEvents enables or disables change event notifications.
func RegisterChangeEvents(c *client.Client, enable bool) error {
	return c.RegisterForChangeEvents(enable)
}

// SetPrettyJSON enables or disables pretty JSON on the speaker.
func SetPrettyJSON(c *client.Client, enable bool) error {
	return c.SetPrettyJSON(enable)
}
