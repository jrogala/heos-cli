// Package cmdutil provides shared helpers for CLI commands.
package cmdutil

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/jrogala/heos-cli/client"
	"github.com/jrogala/heos-cli/config"
	"github.com/spf13/cobra"
)

// NewClient creates a connected HEOS client and returns it with a cleanup function.
func NewClient() (*client.Client, func(), error) {
	cfg := config.Get()
	if cfg.Host == "" {
		return nil, nil, fmt.Errorf("HEOS speaker host not configured. Run 'heos setup' or set HEOS_HOST")
	}
	c := client.New(cfg.Host, cfg.Port)
	if err := c.Connect(); err != nil {
		return nil, nil, err
	}
	return c, func() { c.Close() }, nil
}

// PrintJSON encodes v as indented JSON to stdout.
func PrintJSON(v any) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(v)
}

// IsJSON returns true if the --json persistent flag is set.
func IsJSON(cmd *cobra.Command) bool {
	v, _ := cmd.Root().PersistentFlags().GetBool("json")
	return v
}

// Render outputs data as JSON if --json is set, otherwise calls tableFunc.
func Render(cmd *cobra.Command, data any, tableFunc func()) {
	if IsJSON(cmd) {
		PrintJSON(data)
		return
	}
	tableFunc()
}

// NewTabWriter creates a standard tabwriter for table output.
func NewTabWriter() *tabwriter.Writer {
	return tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
}

// ExitErr prints an error and exits.
func ExitErr(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}
