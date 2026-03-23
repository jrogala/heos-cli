package cmd

import (
	"fmt"

	"github.com/jrogala/heos-cli/internal/cmdutil"
	"github.com/jrogala/heos-cli/pkg/ops"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(browseCmd)
	browseCmd.AddCommand(
		browseSourcesCmd,
		browseSourceInfoCmd,
		browseBrowseCmd,
		browseSearchCmd,
		browseGetSearchCriteriaCmd,
		browsePlayStationCmd,
		browsePlayPresetCmd,
		browsePlayInputCmd,
		browsePlayURLCmd,
		browseAddToQueueCmd,
		browseAddTrackToQueueCmd,
		browsePlaylistsCmd,
		browseRenamePlaylistCmd,
		browseDeletePlaylistCmd,
		browseHistoryCmd,
		browseAlbumMetadataCmd,
		browseSetServiceOptionCmd,
	)

	browseSourceInfoCmd.Flags().String("sid", "", "Source ID")
	browseSourceInfoCmd.MarkFlagRequired("sid")

	browseBrowseCmd.Flags().String("sid", "", "Source ID")
	browseBrowseCmd.MarkFlagRequired("sid")
	browseBrowseCmd.Flags().String("cid", "", "Container ID (optional)")
	browseBrowseCmd.Flags().String("range-start", "", "Range start index")
	browseBrowseCmd.Flags().String("range-end", "", "Range end index")

	browseSearchCmd.Flags().String("sid", "", "Source ID")
	browseSearchCmd.MarkFlagRequired("sid")
	browseSearchCmd.Flags().String("search", "", "Search string")
	browseSearchCmd.MarkFlagRequired("search")
	browseSearchCmd.Flags().String("scid", "", "Search criteria ID")
	browseSearchCmd.MarkFlagRequired("scid")
	browseSearchCmd.Flags().String("range-start", "", "Range start index")
	browseSearchCmd.Flags().String("range-end", "", "Range end index")

	browseGetSearchCriteriaCmd.Flags().String("sid", "", "Source ID")
	browseGetSearchCriteriaCmd.MarkFlagRequired("sid")

	browsePlayStationCmd.Flags().String("pid", "", "Player ID")
	browsePlayStationCmd.MarkFlagRequired("pid")
	browsePlayStationCmd.Flags().String("sid", "", "Source ID")
	browsePlayStationCmd.MarkFlagRequired("sid")
	browsePlayStationCmd.Flags().String("cid", "", "Container ID")
	browsePlayStationCmd.MarkFlagRequired("cid")
	browsePlayStationCmd.Flags().String("mid", "", "Media ID")
	browsePlayStationCmd.MarkFlagRequired("mid")
	browsePlayStationCmd.Flags().String("name", "", "Station name")
	browsePlayStationCmd.MarkFlagRequired("name")

	browsePlayPresetCmd.Flags().String("pid", "", "Player ID")
	browsePlayPresetCmd.MarkFlagRequired("pid")
	browsePlayPresetCmd.Flags().String("preset", "", "Preset position (1+)")
	browsePlayPresetCmd.MarkFlagRequired("preset")

	browsePlayInputCmd.Flags().String("pid", "", "Player ID")
	browsePlayInputCmd.MarkFlagRequired("pid")
	browsePlayInputCmd.Flags().String("input", "", "Input source name")
	browsePlayInputCmd.MarkFlagRequired("input")
	browsePlayInputCmd.Flags().String("spid", "", "Source player ID (for cross-speaker)")

	browsePlayURLCmd.Flags().String("pid", "", "Player ID")
	browsePlayURLCmd.MarkFlagRequired("pid")
	browsePlayURLCmd.Flags().String("url", "", "URL to stream")
	browsePlayURLCmd.MarkFlagRequired("url")

	browseAddToQueueCmd.Flags().String("pid", "", "Player ID")
	browseAddToQueueCmd.MarkFlagRequired("pid")
	browseAddToQueueCmd.Flags().String("sid", "", "Source ID")
	browseAddToQueueCmd.MarkFlagRequired("sid")
	browseAddToQueueCmd.Flags().String("cid", "", "Container ID")
	browseAddToQueueCmd.MarkFlagRequired("cid")
	browseAddToQueueCmd.Flags().String("aid", "1", "Add criteria (1=now, 2=next, 3=end, 4=replace)")

	browseAddTrackToQueueCmd.Flags().String("pid", "", "Player ID")
	browseAddTrackToQueueCmd.MarkFlagRequired("pid")
	browseAddTrackToQueueCmd.Flags().String("sid", "", "Source ID")
	browseAddTrackToQueueCmd.MarkFlagRequired("sid")
	browseAddTrackToQueueCmd.Flags().String("cid", "", "Container ID")
	browseAddTrackToQueueCmd.MarkFlagRequired("cid")
	browseAddTrackToQueueCmd.Flags().String("mid", "", "Media ID")
	browseAddTrackToQueueCmd.MarkFlagRequired("mid")
	browseAddTrackToQueueCmd.Flags().String("aid", "1", "Add criteria (1=now, 2=next, 3=end, 4=replace)")

	browseRenamePlaylistCmd.Flags().String("sid", "", "Source ID")
	browseRenamePlaylistCmd.MarkFlagRequired("sid")
	browseRenamePlaylistCmd.Flags().String("cid", "", "Container ID")
	browseRenamePlaylistCmd.MarkFlagRequired("cid")
	browseRenamePlaylistCmd.Flags().String("name", "", "New playlist name")
	browseRenamePlaylistCmd.MarkFlagRequired("name")

	browseDeletePlaylistCmd.Flags().String("sid", "", "Source ID")
	browseDeletePlaylistCmd.MarkFlagRequired("sid")
	browseDeletePlaylistCmd.Flags().String("cid", "", "Container ID")
	browseDeletePlaylistCmd.MarkFlagRequired("cid")

	browseAlbumMetadataCmd.Flags().String("sid", "", "Source ID")
	browseAlbumMetadataCmd.MarkFlagRequired("sid")
	browseAlbumMetadataCmd.Flags().String("cid", "", "Album ID")
	browseAlbumMetadataCmd.MarkFlagRequired("cid")

	browseSetServiceOptionCmd.Flags().String("sid", "", "Source ID")
	browseSetServiceOptionCmd.MarkFlagRequired("sid")
	browseSetServiceOptionCmd.Flags().String("option", "", "Option ID")
	browseSetServiceOptionCmd.MarkFlagRequired("option")
	browseSetServiceOptionCmd.Flags().String("pid", "", "Player ID")
	browseSetServiceOptionCmd.Flags().String("mid", "", "Media ID")
	browseSetServiceOptionCmd.Flags().String("cid", "", "Container ID")
	browseSetServiceOptionCmd.Flags().String("name", "", "Name")
	browseSetServiceOptionCmd.Flags().String("scid", "", "Search criteria ID")
	browseSetServiceOptionCmd.Flags().String("range-start", "", "Range start")
	browseSetServiceOptionCmd.Flags().String("range-end", "", "Range end")
}

var browseCmd = &cobra.Command{
	Use:     "browse",
	Aliases: []string{"b"},
	Short:   "Browse commands",
}

var browseSourcesCmd = &cobra.Command{
	Use:     "sources",
	Aliases: []string{"src"},
	Short:   "Get music sources",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		sources, err := ops.GetMusicSources(c)
		if err != nil {
			return err
		}
		cmdutil.Render(cmd, sources, func() {
			w := cmdutil.NewTabWriter()
			fmt.Fprintln(w, "SID\tNAME\tTYPE\tAVAILABLE")
			for _, s := range sources {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", s.SID, s.Name, s.Type, s.Available)
			}
			w.Flush()
		})
		return nil
	},
}

var browseSourceInfoCmd = &cobra.Command{
	Use:   "source-info",
	Short: "Get source info",
	RunE: func(cmd *cobra.Command, args []string) error {
		sid, _ := cmd.Flags().GetString("sid")
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		sources, err := ops.GetSourceInfo(c, sid)
		if err != nil {
			return err
		}
		cmdutil.Render(cmd, sources, func() {
			for _, s := range sources {
				fmt.Printf("Name:      %s\n", s.Name)
				fmt.Printf("SID:       %s\n", s.SID)
				fmt.Printf("Type:      %s\n", s.Type)
				fmt.Printf("Available: %s\n", s.Available)
			}
		})
		return nil
	},
}

var browseBrowseCmd = &cobra.Command{
	Use:   "browse",
	Short: "Browse a source or container",
	RunE: func(cmd *cobra.Command, args []string) error {
		sid, _ := cmd.Flags().GetString("sid")
		cid, _ := cmd.Flags().GetString("cid")
		rangeStart, _ := cmd.Flags().GetString("range-start")
		rangeEnd, _ := cmd.Flags().GetString("range-end")

		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()

		var items []ops.BrowseItem
		if cid != "" {
			items, err = ops.BrowseSourceContainers(c, sid, cid, rangeStart, rangeEnd)
		} else {
			items, err = ops.BrowseSource(c, sid, rangeStart, rangeEnd)
		}
		if err != nil {
			return err
		}

		cmdutil.Render(cmd, items, func() {
			w := cmdutil.NewTabWriter()
			fmt.Fprintln(w, "TYPE\tNAME\tCID/MID")
			for _, item := range items {
				id := item.CID
				if id == "" {
					id = item.MID
				}
				fmt.Fprintf(w, "%s\t%s\t%s\n", item.Type, item.Name, id)
			}
			w.Flush()
		})
		return nil
	},
}

var browseSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search a music source",
	RunE: func(cmd *cobra.Command, args []string) error {
		sid, _ := cmd.Flags().GetString("sid")
		search, _ := cmd.Flags().GetString("search")
		scid, _ := cmd.Flags().GetString("scid")
		rangeStart, _ := cmd.Flags().GetString("range-start")
		rangeEnd, _ := cmd.Flags().GetString("range-end")

		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()

		items, err := ops.Search(c, sid, search, scid, rangeStart, rangeEnd)
		if err != nil {
			return err
		}

		cmdutil.Render(cmd, items, func() {
			w := cmdutil.NewTabWriter()
			fmt.Fprintln(w, "TYPE\tNAME\tCID/MID")
			for _, item := range items {
				id := item.CID
				if id == "" {
					id = item.MID
				}
				fmt.Fprintf(w, "%s\t%s\t%s\n", item.Type, item.Name, id)
			}
			w.Flush()
		})
		return nil
	},
}

var browseGetSearchCriteriaCmd = &cobra.Command{
	Use:   "get-search-criteria",
	Short: "Get search criteria for a source",
	RunE: func(cmd *cobra.Command, args []string) error {
		sid, _ := cmd.Flags().GetString("sid")
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		criteria, err := ops.GetSearchCriteria(c, sid)
		if err != nil {
			return err
		}
		cmdutil.Render(cmd, criteria, func() {
			w := cmdutil.NewTabWriter()
			fmt.Fprintln(w, "SCID\tNAME")
			for _, sc := range criteria {
				fmt.Fprintf(w, "%s\t%s\n", sc.SCID, sc.Name)
			}
			w.Flush()
		})
		return nil
	},
}

var browsePlayStationCmd = &cobra.Command{
	Use:   "play-station",
	Short: "Play a station",
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, _ := cmd.Flags().GetString("pid")
		sid, _ := cmd.Flags().GetString("sid")
		cid, _ := cmd.Flags().GetString("cid")
		mid, _ := cmd.Flags().GetString("mid")
		name, _ := cmd.Flags().GetString("name")

		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		if err := ops.PlayStation(c, pid, sid, cid, mid, name); err != nil {
			return err
		}
		fmt.Printf("Playing station: %s\n", name)
		return nil
	},
}

var browsePlayPresetCmd = &cobra.Command{
	Use:   "play-preset",
	Short: "Play a preset station from HEOS Favorites",
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, _ := cmd.Flags().GetString("pid")
		preset, _ := cmd.Flags().GetString("preset")

		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		if err := ops.PlayPreset(c, pid, preset); err != nil {
			return err
		}
		fmt.Printf("Playing preset %s\n", preset)
		return nil
	},
}

var browsePlayInputCmd = &cobra.Command{
	Use:   "play-input",
	Short: "Play an input source",
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, _ := cmd.Flags().GetString("pid")
		input, _ := cmd.Flags().GetString("input")
		spid, _ := cmd.Flags().GetString("spid")

		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		if err := ops.PlayInput(c, pid, input, spid); err != nil {
			return err
		}
		fmt.Printf("Playing input: %s\n", input)
		return nil
	},
}

var browsePlayURLCmd = &cobra.Command{
	Use:   "play-url",
	Short: "Play a URL stream",
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, _ := cmd.Flags().GetString("pid")
		url, _ := cmd.Flags().GetString("url")

		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		if err := ops.PlayURL(c, pid, url); err != nil {
			return err
		}
		fmt.Println("Playing URL")
		return nil
	},
}

var browseAddToQueueCmd = &cobra.Command{
	Use:   "add-to-queue",
	Short: "Add container to queue",
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, _ := cmd.Flags().GetString("pid")
		sid, _ := cmd.Flags().GetString("sid")
		cid, _ := cmd.Flags().GetString("cid")
		aid, _ := cmd.Flags().GetString("aid")

		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		if err := ops.AddToQueue(c, pid, sid, cid, aid); err != nil {
			return err
		}
		fmt.Println("Added to queue")
		return nil
	},
}

var browseAddTrackToQueueCmd = &cobra.Command{
	Use:   "add-track-to-queue",
	Short: "Add track to queue",
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, _ := cmd.Flags().GetString("pid")
		sid, _ := cmd.Flags().GetString("sid")
		cid, _ := cmd.Flags().GetString("cid")
		mid, _ := cmd.Flags().GetString("mid")
		aid, _ := cmd.Flags().GetString("aid")

		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		if err := ops.AddTrackToQueue(c, pid, sid, cid, mid, aid); err != nil {
			return err
		}
		fmt.Println("Track added to queue")
		return nil
	},
}

var browsePlaylistsCmd = &cobra.Command{
	Use:   "playlists",
	Short: "Get HEOS playlists (browses SID 1025)",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		items, err := ops.BrowseSource(c, "1025", "", "")
		if err != nil {
			return err
		}
		cmdutil.Render(cmd, items, func() {
			w := cmdutil.NewTabWriter()
			fmt.Fprintln(w, "CID\tNAME\tTYPE")
			for _, item := range items {
				fmt.Fprintf(w, "%s\t%s\t%s\n", item.CID, item.Name, item.Type)
			}
			w.Flush()
		})
		return nil
	},
}

var browseRenamePlaylistCmd = &cobra.Command{
	Use:   "rename-playlist",
	Short: "Rename a HEOS playlist",
	RunE: func(cmd *cobra.Command, args []string) error {
		sid, _ := cmd.Flags().GetString("sid")
		cid, _ := cmd.Flags().GetString("cid")
		name, _ := cmd.Flags().GetString("name")

		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		if err := ops.RenamePlaylist(c, sid, cid, name); err != nil {
			return err
		}
		fmt.Printf("Playlist renamed to '%s'\n", name)
		return nil
	},
}

var browseDeletePlaylistCmd = &cobra.Command{
	Use:   "delete-playlist",
	Short: "Delete a HEOS playlist",
	RunE: func(cmd *cobra.Command, args []string) error {
		sid, _ := cmd.Flags().GetString("sid")
		cid, _ := cmd.Flags().GetString("cid")

		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		if err := ops.DeletePlaylist(c, sid, cid); err != nil {
			return err
		}
		fmt.Println("Playlist deleted")
		return nil
	},
}

var browseHistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "Get HEOS history (browses SID 1026)",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		items, err := ops.BrowseSource(c, "1026", "", "")
		if err != nil {
			return err
		}
		cmdutil.Render(cmd, items, func() {
			w := cmdutil.NewTabWriter()
			fmt.Fprintln(w, "TYPE\tNAME\tMID")
			for _, item := range items {
				fmt.Fprintf(w, "%s\t%s\t%s\n", item.Type, item.Name, item.MID)
			}
			w.Flush()
		})
		return nil
	},
}

var browseAlbumMetadataCmd = &cobra.Command{
	Use:   "album-metadata",
	Short: "Retrieve album metadata",
	RunE: func(cmd *cobra.Command, args []string) error {
		sid, _ := cmd.Flags().GetString("sid")
		cid, _ := cmd.Flags().GetString("cid")

		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		meta, err := c.RetrieveMetadata(sid, cid)
		if err != nil {
			return err
		}
		cmdutil.Render(cmd, meta, func() {
			for _, m := range meta {
				fmt.Printf("Album ID: %s\n", m.AlbumID)
				for _, img := range m.Images {
					fmt.Printf("  Image: %s (width: %s)\n", img.ImageURL, img.Width)
				}
			}
		})
		return nil
	},
}

var browseSetServiceOptionCmd = &cobra.Command{
	Use:   "set-service-option",
	Short: "Set a service option",
	RunE: func(cmd *cobra.Command, args []string) error {
		sid, _ := cmd.Flags().GetString("sid")
		option, _ := cmd.Flags().GetString("option")

		optArgs := map[string]string{"sid": sid, "option": option}

		if v, _ := cmd.Flags().GetString("pid"); v != "" {
			optArgs["pid"] = v
		}
		if v, _ := cmd.Flags().GetString("mid"); v != "" {
			optArgs["mid"] = v
		}
		if v, _ := cmd.Flags().GetString("cid"); v != "" {
			optArgs["cid"] = v
		}
		if v, _ := cmd.Flags().GetString("name"); v != "" {
			optArgs["name"] = v
		}
		if v, _ := cmd.Flags().GetString("scid"); v != "" {
			optArgs["scid"] = v
		}
		if v, _ := cmd.Flags().GetString("range-start"); v != "" {
			re, _ := cmd.Flags().GetString("range-end")
			optArgs["range"] = v + "," + re
		}

		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		if err := ops.SetServiceOption(c, optArgs); err != nil {
			return err
		}
		fmt.Println("Service option set")
		return nil
	},
}
