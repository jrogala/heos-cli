package cmd

import (
	"os"

	"github.com/jrogala/heos-cli/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "heos",
	Short: "CLI for controlling Denon HEOS speakers",
	Long: `A command-line interface for controlling Denon HEOS speakers via the HEOS CLI Protocol.

Quick examples:
  heos player list                                  List all players (get PIDs)
  heos player get-volume --pid PID                  Get player volume
  heos player set-volume --pid PID --level 20       Set volume (0-100)
  heos player now-playing --pid PID                 Show current track
  heos player play --pid PID                        Resume playback
  heos player pause --pid PID                       Pause playback
  heos player mute --pid PID                        Mute player
  heos player unmute --pid PID                      Unmute player
  heos group list                                   List groups
  heos group set-volume --gid GID --level 20        Set group volume
  heos browse sources                               List music sources
  heos browse play-url --pid PID --url URL          Play audio URL on speaker
  heos system check-account                         Check HEOS account
  heos system heart-beat                            Test connection`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(config.Init)
	rootCmd.PersistentFlags().Bool("json", false, "Output raw JSON responses")
}
