package client

import "fmt"

// HEOS error codes from protocol spec section 6.2
const (
	ErrUnrecognizedCommand   = 1
	ErrInvalidID             = 2
	ErrWrongArguments        = 3
	ErrRequestedDataNA       = 4
	ErrResourceNA            = 5
	ErrInvalidCredentials    = 6
	ErrCommandNotExecuted    = 7
	ErrUserNotLoggedIn       = 8
	ErrParameterOutOfRange   = 9
	ErrUserNotFound          = 10
	ErrInternalError         = 11
	ErrSystemError           = 12
	ErrProcessingPrevious    = 13
	ErrMediaCantBePlayed     = 14
	ErrOptionNotSupported    = 15
	ErrTooManyCommands       = 16
	ErrReachedSkipLimit      = 17
)

// HEOSError represents an error returned by the HEOS system.
type HEOSError struct {
	EID     int
	Text    string
	SysErno int // only for eid=12
}

func (e *HEOSError) Error() string {
	if e.SysErno != 0 {
		return fmt.Sprintf("heos error %d: %s (syserrno=%d)", e.EID, e.Text, e.SysErno)
	}
	return fmt.Sprintf("heos error %d: %s", e.EID, e.Text)
}

// parseHEOSError creates a HEOSError from a parsed message map.
func parseHEOSError(msg map[string]string) *HEOSError {
	eid := 0
	fmt.Sscanf(msg["eid"], "%d", &eid)

	syserrno := 0
	if v, ok := msg["syserrno"]; ok {
		fmt.Sscanf(v, "%d", &syserrno)
	}

	return &HEOSError{
		EID:     eid,
		Text:    msg["text"],
		SysErno: syserrno,
	}
}
