package cmd

import (
	"fmt"

	"github.com/jrogala/heos-cli/internal/cmdutil"
	"github.com/jrogala/heos-cli/pkg/ops"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(playerCmd)
	playerCmd.AddCommand(
		playerListCmd,
		playerInfoCmd,
		playerGetStateCmd,
		playerPlayCmd,
		playerPauseCmd,
		playerStopCmd,
		playerNowPlayingCmd,
		playerGetVolumeCmd,
		playerSetVolumeCmd,
		playerVolumeUpCmd,
		playerVolumeDownCmd,
		playerGetMuteCmd,
		playerMuteCmd,
		playerUnmuteCmd,
		playerToggleMuteCmd,
		playerGetPlayModeCmd,
		playerSetPlayModeCmd,
		playerQueueCmd,
		playerPlayQueueCmd,
		playerRemoveFromQueueCmd,
		playerSaveQueueCmd,
		playerClearQueueCmd,
		playerMoveQueueCmd,
		playerNextCmd,
		playerPreviousCmd,
		playerSetQuickSelectCmd,
		playerPlayQuickSelectCmd,
		playerGetQuickSelectsCmd,
		playerCheckUpdateCmd,
	)

	// --pid flag for all player subcommands that need it
	pidCmds := []*cobra.Command{
		playerInfoCmd, playerGetStateCmd,
		playerPlayCmd, playerPauseCmd, playerStopCmd,
		playerNowPlayingCmd,
		playerGetVolumeCmd, playerSetVolumeCmd,
		playerVolumeUpCmd, playerVolumeDownCmd,
		playerGetMuteCmd, playerMuteCmd, playerUnmuteCmd, playerToggleMuteCmd,
		playerGetPlayModeCmd, playerSetPlayModeCmd,
		playerQueueCmd, playerPlayQueueCmd,
		playerRemoveFromQueueCmd, playerSaveQueueCmd,
		playerClearQueueCmd, playerMoveQueueCmd,
		playerNextCmd, playerPreviousCmd,
		playerSetQuickSelectCmd, playerPlayQuickSelectCmd,
		playerGetQuickSelectsCmd,
		playerCheckUpdateCmd,
	}
	for _, cmd := range pidCmds {
		cmd.Flags().String("pid", "", "Player ID")
		cmd.MarkFlagRequired("pid")
	}

	playerSetVolumeCmd.Flags().Int("level", 0, "Volume level (0-100)")
	playerSetVolumeCmd.MarkFlagRequired("level")

	playerVolumeUpCmd.Flags().Int("step", 5, "Volume step (1-10)")
	playerVolumeDownCmd.Flags().Int("step", 5, "Volume step (1-10)")

	playerSetPlayModeCmd.Flags().String("repeat", "off", "Repeat mode (on_all, on_one, off)")
	playerSetPlayModeCmd.Flags().String("shuffle", "off", "Shuffle mode (on, off)")

	playerQueueCmd.Flags().String("range-start", "", "Range start index")
	playerQueueCmd.Flags().String("range-end", "", "Range end index")

	playerPlayQueueCmd.Flags().String("qid", "", "Queue item ID")
	playerPlayQueueCmd.MarkFlagRequired("qid")

	playerRemoveFromQueueCmd.Flags().String("qid", "", "Queue item ID(s), comma-separated")
	playerRemoveFromQueueCmd.MarkFlagRequired("qid")

	playerSaveQueueCmd.Flags().String("name", "", "Playlist name")
	playerSaveQueueCmd.MarkFlagRequired("name")

	playerMoveQueueCmd.Flags().String("sqid", "", "Source queue ID(s), comma-separated")
	playerMoveQueueCmd.Flags().String("dqid", "", "Destination queue ID")
	playerMoveQueueCmd.MarkFlagRequired("sqid")
	playerMoveQueueCmd.MarkFlagRequired("dqid")

	playerSetQuickSelectCmd.Flags().String("id", "", "QuickSelect ID (1-6)")
	playerSetQuickSelectCmd.MarkFlagRequired("id")
	playerPlayQuickSelectCmd.Flags().String("id", "", "QuickSelect ID (1-6)")
	playerPlayQuickSelectCmd.MarkFlagRequired("id")
	playerGetQuickSelectsCmd.Flags().String("id", "", "QuickSelect ID (1-6, optional)")
}

var playerCmd = &cobra.Command{
	Use:     "player",
	Aliases: []string{"p"},
	Short:   "Player commands",
}

var playerListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List available players",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		players, err := ops.ListPlayers(c)
		if err != nil {
			return err
		}
		cmdutil.Render(cmd, players, func() {
			w := cmdutil.NewTabWriter()
			fmt.Fprintln(w, "PID\tNAME\tMODEL\tNETWORK\tSERIAL")
			for _, p := range players {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", p.PID, p.Name, p.Model, p.Network, p.Serial)
			}
			w.Flush()
		})
		return nil
	},
}

var playerInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get player info",
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, _ := cmd.Flags().GetString("pid")
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		player, err := ops.GetPlayerInfo(c, pid)
		if err != nil {
			return err
		}
		cmdutil.Render(cmd, player, func() {
			fmt.Printf("Name:    %s\n", player.Name)
			fmt.Printf("PID:     %s\n", player.PID)
			fmt.Printf("Model:   %s\n", player.Model)
			fmt.Printf("Version: %s\n", player.Version)
			fmt.Printf("Network: %s\n", player.Network)
			if player.Serial != "" {
				fmt.Printf("Serial:  %s\n", player.Serial)
			}
		})
		return nil
	},
}

var playerGetStateCmd = &cobra.Command{
	Use:   "get-state",
	Short: "Get player play state",
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, _ := cmd.Flags().GetString("pid")
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		state, err := ops.GetPlayState(c, pid)
		if err != nil {
			return err
		}
		cmdutil.Render(cmd, map[string]string{"state": state}, func() {
			fmt.Println(state)
		})
		return nil
	},
}

var playerPlayCmd = &cobra.Command{
	Use:   "play",
	Short: "Start playback",
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, _ := cmd.Flags().GetString("pid")
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		if err := ops.SetPlayState(c, pid, "play"); err != nil {
			return err
		}
		fmt.Println("Playing")
		return nil
	},
}

var playerPauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "Pause playback",
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, _ := cmd.Flags().GetString("pid")
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		if err := ops.SetPlayState(c, pid, "pause"); err != nil {
			return err
		}
		fmt.Println("Paused")
		return nil
	},
}

var playerStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop playback",
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, _ := cmd.Flags().GetString("pid")
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		if err := ops.SetPlayState(c, pid, "stop"); err != nil {
			return err
		}
		fmt.Println("Stopped")
		return nil
	},
}

var playerNowPlayingCmd = &cobra.Command{
	Use:     "now-playing",
	Aliases: []string{"np"},
	Short:   "Get now playing media",
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, _ := cmd.Flags().GetString("pid")
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		np, err := ops.GetNowPlaying(c, pid)
		if err != nil {
			return err
		}
		cmdutil.Render(cmd, np, func() {
			fmt.Printf("Type:    %s\n", np.Type)
			fmt.Printf("Song:    %s\n", np.Song)
			fmt.Printf("Artist:  %s\n", np.Artist)
			fmt.Printf("Album:   %s\n", np.Album)
			if np.Station != "" {
				fmt.Printf("Station: %s\n", np.Station)
			}
		})
		return nil
	},
}

var playerGetVolumeCmd = &cobra.Command{
	Use:   "get-volume",
	Short: "Get player volume",
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, _ := cmd.Flags().GetString("pid")
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		level, err := ops.GetVolume(c, pid)
		if err != nil {
			return err
		}
		cmdutil.Render(cmd, map[string]int{"level": level}, func() {
			fmt.Println(level)
		})
		return nil
	},
}

var playerSetVolumeCmd = &cobra.Command{
	Use:   "set-volume",
	Short: "Set player volume (0-100)",
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, _ := cmd.Flags().GetString("pid")
		level, _ := cmd.Flags().GetInt("level")
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		if err := ops.SetVolume(c, pid, level); err != nil {
			return err
		}
		fmt.Printf("Volume set to %d\n", level)
		return nil
	},
}

var playerVolumeUpCmd = &cobra.Command{
	Use:   "volume-up",
	Short: "Increase player volume",
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, _ := cmd.Flags().GetString("pid")
		step, _ := cmd.Flags().GetInt("step")
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		if err := ops.VolumeUp(c, pid, step); err != nil {
			return err
		}
		fmt.Println("Volume increased")
		return nil
	},
}

var playerVolumeDownCmd = &cobra.Command{
	Use:   "volume-down",
	Short: "Decrease player volume",
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, _ := cmd.Flags().GetString("pid")
		step, _ := cmd.Flags().GetInt("step")
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		if err := ops.VolumeDown(c, pid, step); err != nil {
			return err
		}
		fmt.Println("Volume decreased")
		return nil
	},
}

var playerGetMuteCmd = &cobra.Command{
	Use:   "get-mute",
	Short: "Get player mute state",
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, _ := cmd.Flags().GetString("pid")
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		state, err := ops.GetMute(c, pid)
		if err != nil {
			return err
		}
		cmdutil.Render(cmd, map[string]string{"state": state}, func() {
			fmt.Println(state)
		})
		return nil
	},
}

var playerMuteCmd = &cobra.Command{
	Use:   "mute",
	Short: "Mute player",
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, _ := cmd.Flags().GetString("pid")
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		if err := ops.SetMute(c, pid, "on"); err != nil {
			return err
		}
		fmt.Println("Muted")
		return nil
	},
}

var playerUnmuteCmd = &cobra.Command{
	Use:   "unmute",
	Short: "Unmute player",
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, _ := cmd.Flags().GetString("pid")
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		if err := ops.SetMute(c, pid, "off"); err != nil {
			return err
		}
		fmt.Println("Unmuted")
		return nil
	},
}

var playerToggleMuteCmd = &cobra.Command{
	Use:   "toggle-mute",
	Short: "Toggle player mute",
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, _ := cmd.Flags().GetString("pid")
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		if err := ops.ToggleMute(c, pid); err != nil {
			return err
		}
		fmt.Println("Mute toggled")
		return nil
	},
}

var playerGetPlayModeCmd = &cobra.Command{
	Use:   "get-play-mode",
	Short: "Get play mode (repeat/shuffle)",
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, _ := cmd.Flags().GetString("pid")
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		mode, err := ops.GetPlayMode(c, pid)
		if err != nil {
			return err
		}
		cmdutil.Render(cmd, mode, func() {
			fmt.Printf("Repeat:  %s\n", mode.Repeat)
			fmt.Printf("Shuffle: %s\n", mode.Shuffle)
		})
		return nil
	},
}

var playerSetPlayModeCmd = &cobra.Command{
	Use:   "set-play-mode",
	Short: "Set play mode (repeat/shuffle)",
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, _ := cmd.Flags().GetString("pid")
		repeat, _ := cmd.Flags().GetString("repeat")
		shuffle, _ := cmd.Flags().GetString("shuffle")
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		if err := ops.SetPlayMode(c, pid, repeat, shuffle); err != nil {
			return err
		}
		fmt.Printf("Play mode set: repeat=%s, shuffle=%s\n", repeat, shuffle)
		return nil
	},
}

var playerQueueCmd = &cobra.Command{
	Use:     "queue",
	Aliases: []string{"q"},
	Short:   "Get player queue",
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, _ := cmd.Flags().GetString("pid")
		rangeStart, _ := cmd.Flags().GetString("range-start")
		rangeEnd, _ := cmd.Flags().GetString("range-end")
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		items, err := ops.GetQueue(c, pid, rangeStart, rangeEnd)
		if err != nil {
			return err
		}
		cmdutil.Render(cmd, items, func() {
			w := cmdutil.NewTabWriter()
			fmt.Fprintln(w, "QID\tSONG\tARTIST\tALBUM")
			for _, item := range items {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", item.QID, item.Song, item.Artist, item.Album)
			}
			w.Flush()
		})
		return nil
	},
}

var playerPlayQueueCmd = &cobra.Command{
	Use:   "play-queue",
	Short: "Play a queue item",
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, _ := cmd.Flags().GetString("pid")
		qid, _ := cmd.Flags().GetString("qid")
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		if err := ops.PlayQueueItem(c, pid, qid); err != nil {
			return err
		}
		fmt.Println("Playing queue item")
		return nil
	},
}

var playerRemoveFromQueueCmd = &cobra.Command{
	Use:   "remove-from-queue",
	Short: "Remove item(s) from queue",
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, _ := cmd.Flags().GetString("pid")
		qid, _ := cmd.Flags().GetString("qid")
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		if err := ops.RemoveFromQueue(c, pid, qid); err != nil {
			return err
		}
		fmt.Println("Removed from queue")
		return nil
	},
}

var playerSaveQueueCmd = &cobra.Command{
	Use:   "save-queue",
	Short: "Save queue as playlist",
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, _ := cmd.Flags().GetString("pid")
		name, _ := cmd.Flags().GetString("name")
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		if err := ops.SaveQueue(c, pid, name); err != nil {
			return err
		}
		fmt.Printf("Queue saved as '%s'\n", name)
		return nil
	},
}

var playerClearQueueCmd = &cobra.Command{
	Use:   "clear-queue",
	Short: "Clear the queue",
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, _ := cmd.Flags().GetString("pid")
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		if err := ops.ClearQueue(c, pid); err != nil {
			return err
		}
		fmt.Println("Queue cleared")
		return nil
	},
}

var playerMoveQueueCmd = &cobra.Command{
	Use:   "move-queue",
	Short: "Move queue items",
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, _ := cmd.Flags().GetString("pid")
		sqid, _ := cmd.Flags().GetString("sqid")
		dqid, _ := cmd.Flags().GetString("dqid")
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		if err := ops.MoveQueue(c, pid, sqid, dqid); err != nil {
			return err
		}
		fmt.Println("Queue items moved")
		return nil
	},
}

var playerNextCmd = &cobra.Command{
	Use:   "next",
	Short: "Play next track",
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, _ := cmd.Flags().GetString("pid")
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		if err := ops.PlayNext(c, pid); err != nil {
			return err
		}
		fmt.Println("Next")
		return nil
	},
}

var playerPreviousCmd = &cobra.Command{
	Use:     "previous",
	Aliases: []string{"prev"},
	Short:   "Play previous track",
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, _ := cmd.Flags().GetString("pid")
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		if err := ops.PlayPrevious(c, pid); err != nil {
			return err
		}
		fmt.Println("Previous")
		return nil
	},
}

var playerSetQuickSelectCmd = &cobra.Command{
	Use:   "set-quickselect",
	Short: "Set QuickSelect (LS AVR only)",
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, _ := cmd.Flags().GetString("pid")
		id, _ := cmd.Flags().GetString("id")
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		if err := ops.SetQuickSelect(c, pid, id); err != nil {
			return err
		}
		fmt.Printf("QuickSelect %s set\n", id)
		return nil
	},
}

var playerPlayQuickSelectCmd = &cobra.Command{
	Use:   "play-quickselect",
	Short: "Play QuickSelect (LS AVR only)",
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, _ := cmd.Flags().GetString("pid")
		id, _ := cmd.Flags().GetString("id")
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		if err := ops.PlayQuickSelect(c, pid, id); err != nil {
			return err
		}
		fmt.Printf("Playing QuickSelect %s\n", id)
		return nil
	},
}

var playerGetQuickSelectsCmd = &cobra.Command{
	Use:   "get-quickselects",
	Short: "Get QuickSelects (LS AVR only)",
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, _ := cmd.Flags().GetString("pid")
		id, _ := cmd.Flags().GetString("id")
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		qs, err := ops.GetQuickSelects(c, pid, id)
		if err != nil {
			return err
		}
		cmdutil.Render(cmd, qs, func() {
			w := cmdutil.NewTabWriter()
			fmt.Fprintln(w, "ID\tNAME")
			for _, q := range qs {
				fmt.Fprintf(w, "%d\t%s\n", q.ID, q.Name)
			}
			w.Flush()
		})
		return nil
	},
}

var playerCheckUpdateCmd = &cobra.Command{
	Use:   "check-update",
	Short: "Check for firmware update",
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, _ := cmd.Flags().GetString("pid")
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		update, err := ops.CheckFirmwareUpdate(c, pid)
		if err != nil {
			return err
		}
		cmdutil.Render(cmd, map[string]string{"update": update}, func() {
			fmt.Println(update)
		})
		return nil
	},
}
