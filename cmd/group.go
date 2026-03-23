package cmd

import (
	"fmt"

	"github.com/jrogala/heos-cli/internal/cmdutil"
	"github.com/jrogala/heos-cli/pkg/ops"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(groupCmd)
	groupCmd.AddCommand(
		groupListCmd,
		groupInfoCmd,
		groupSetCmd,
		groupGetVolumeCmd,
		groupSetVolumeCmd,
		groupVolumeUpCmd,
		groupVolumeDownCmd,
		groupGetMuteCmd,
		groupSetMuteCmd,
		groupToggleMuteCmd,
	)

	// --gid flag for commands that need it
	gidCmds := []*cobra.Command{
		groupInfoCmd,
		groupGetVolumeCmd, groupSetVolumeCmd,
		groupVolumeUpCmd, groupVolumeDownCmd,
		groupGetMuteCmd, groupSetMuteCmd, groupToggleMuteCmd,
	}
	for _, cmd := range gidCmds {
		cmd.Flags().String("gid", "", "Group ID")
		cmd.MarkFlagRequired("gid")
	}

	groupSetCmd.Flags().String("pid", "", "Player IDs (comma-separated, first is leader)")
	groupSetCmd.MarkFlagRequired("pid")

	groupSetVolumeCmd.Flags().Int("level", 0, "Volume level (0-100)")
	groupSetVolumeCmd.MarkFlagRequired("level")

	groupVolumeUpCmd.Flags().Int("step", 5, "Volume step (1-10)")
	groupVolumeDownCmd.Flags().Int("step", 5, "Volume step (1-10)")

	groupSetMuteCmd.Flags().String("state", "", "Mute state (on/off)")
	groupSetMuteCmd.MarkFlagRequired("state")
}

var groupCmd = &cobra.Command{
	Use:     "group",
	Aliases: []string{"g"},
	Short:   "Group commands",
}

var groupListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List groups",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		groups, err := ops.ListGroups(c)
		if err != nil {
			return err
		}
		cmdutil.Render(cmd, groups, func() {
			if len(groups) == 0 {
				fmt.Println("No groups")
				return
			}
			w := cmdutil.NewTabWriter()
			fmt.Fprintln(w, "GID\tNAME\tPLAYERS")
			for _, g := range groups {
				fmt.Fprintf(w, "%s\t%s\t%d\n", g.GID, g.Name, len(g.Players))
			}
			w.Flush()
		})
		return nil
	},
}

var groupInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get group info",
	RunE: func(cmd *cobra.Command, args []string) error {
		gid, _ := cmd.Flags().GetString("gid")
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		group, err := ops.GetGroupInfo(c, gid)
		if err != nil {
			return err
		}
		cmdutil.Render(cmd, group, func() {
			fmt.Printf("Name: %s\n", group.Name)
			fmt.Printf("GID:  %s\n", group.GID)
			fmt.Println("Players:")
			for _, p := range group.Players {
				fmt.Printf("  %s (PID: %s, Role: %s)\n", p.Name, p.PID, p.Role)
			}
		})
		return nil
	},
}

var groupSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Create, modify, or ungroup players",
	Long:  "Set group: first PID is leader. Single PID ungroups.",
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, _ := cmd.Flags().GetString("pid")
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		result, err := ops.SetGroup(c, pid)
		if err != nil {
			return err
		}
		cmdutil.Render(cmd, result, func() {
			if result.GID != "" {
				fmt.Printf("Group set (GID: %s)\n", result.GID)
			} else {
				fmt.Println("Group updated")
			}
		})
		return nil
	},
}

var groupGetVolumeCmd = &cobra.Command{
	Use:   "get-volume",
	Short: "Get group volume",
	RunE: func(cmd *cobra.Command, args []string) error {
		gid, _ := cmd.Flags().GetString("gid")
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		level, err := ops.GetGroupVolume(c, gid)
		if err != nil {
			return err
		}
		cmdutil.Render(cmd, map[string]int{"level": level}, func() {
			fmt.Println(level)
		})
		return nil
	},
}

var groupSetVolumeCmd = &cobra.Command{
	Use:   "set-volume",
	Short: "Set group volume (0-100)",
	RunE: func(cmd *cobra.Command, args []string) error {
		gid, _ := cmd.Flags().GetString("gid")
		level, _ := cmd.Flags().GetInt("level")
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		if err := ops.SetGroupVolume(c, gid, level); err != nil {
			return err
		}
		fmt.Printf("Group volume set to %d\n", level)
		return nil
	},
}

var groupVolumeUpCmd = &cobra.Command{
	Use:   "volume-up",
	Short: "Increase group volume",
	RunE: func(cmd *cobra.Command, args []string) error {
		gid, _ := cmd.Flags().GetString("gid")
		step, _ := cmd.Flags().GetInt("step")
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		if err := ops.GroupVolumeUp(c, gid, step); err != nil {
			return err
		}
		fmt.Println("Group volume increased")
		return nil
	},
}

var groupVolumeDownCmd = &cobra.Command{
	Use:   "volume-down",
	Short: "Decrease group volume",
	RunE: func(cmd *cobra.Command, args []string) error {
		gid, _ := cmd.Flags().GetString("gid")
		step, _ := cmd.Flags().GetInt("step")
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		if err := ops.GroupVolumeDown(c, gid, step); err != nil {
			return err
		}
		fmt.Println("Group volume decreased")
		return nil
	},
}

var groupGetMuteCmd = &cobra.Command{
	Use:   "get-mute",
	Short: "Get group mute state",
	RunE: func(cmd *cobra.Command, args []string) error {
		gid, _ := cmd.Flags().GetString("gid")
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		state, err := ops.GetGroupMute(c, gid)
		if err != nil {
			return err
		}
		cmdutil.Render(cmd, map[string]string{"state": state}, func() {
			fmt.Println(state)
		})
		return nil
	},
}

var groupSetMuteCmd = &cobra.Command{
	Use:   "set-mute",
	Short: "Set group mute (on/off)",
	RunE: func(cmd *cobra.Command, args []string) error {
		gid, _ := cmd.Flags().GetString("gid")
		state, _ := cmd.Flags().GetString("state")
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		if err := ops.SetGroupMute(c, gid, state); err != nil {
			return err
		}
		fmt.Printf("Group mute: %s\n", state)
		return nil
	},
}

var groupToggleMuteCmd = &cobra.Command{
	Use:   "toggle-mute",
	Short: "Toggle group mute",
	RunE: func(cmd *cobra.Command, args []string) error {
		gid, _ := cmd.Flags().GetString("gid")
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		if err := ops.ToggleGroupMute(c, gid); err != nil {
			return err
		}
		fmt.Println("Group mute toggled")
		return nil
	},
}
