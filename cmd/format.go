package cmd

import (
	"github.com/jrogala/heos-cli/internal/cmdutil"
)

// Re-export cmdutil helpers for backward compat with setup.go
var exitErr = cmdutil.ExitErr
