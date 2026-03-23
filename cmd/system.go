package cmd

import (
	"fmt"

	"github.com/jrogala/heos-cli/internal/cmdutil"
	"github.com/jrogala/heos-cli/pkg/ops"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(systemCmd)
	systemCmd.AddCommand(
		registerEventsCmd,
		checkAccountCmd,
		signInCmd,
		signOutCmd,
		heartBeatCmd,
		rebootCmd,
		prettifyJSONCmd,
	)

	registerEventsCmd.Flags().Bool("enable", true, "Enable (true) or disable (false) events")
	signInCmd.Flags().String("username", "", "HEOS account username")
	signInCmd.Flags().String("password", "", "HEOS account password")
	prettifyJSONCmd.Flags().Bool("enable", true, "Enable (true) or disable (false) pretty JSON")
}

var systemCmd = &cobra.Command{
	Use:     "system",
	Aliases: []string{"sys"},
	Short:   "System commands",
}

var registerEventsCmd = &cobra.Command{
	Use:   "register-events",
	Short: "Register or unregister for change events",
	RunE: func(cmd *cobra.Command, args []string) error {
		enable, _ := cmd.Flags().GetBool("enable")
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		if err := ops.RegisterChangeEvents(c, enable); err != nil {
			return err
		}
		if enable {
			fmt.Println("Change events enabled")
		} else {
			fmt.Println("Change events disabled")
		}
		return nil
	},
}

var checkAccountCmd = &cobra.Command{
	Use:   "check-account",
	Short: "Check HEOS account status",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		info, err := ops.CheckAccount(c)
		if err != nil {
			return err
		}
		cmdutil.Render(cmd, info, func() {
			if info.SignedIn {
				fmt.Printf("Signed in as: %s\n", info.Username)
			} else {
				fmt.Println("Signed out")
			}
		})
		return nil
	},
}

var signInCmd = &cobra.Command{
	Use:   "sign-in",
	Short: "Sign in to HEOS account",
	RunE: func(cmd *cobra.Command, args []string) error {
		un, _ := cmd.Flags().GetString("username")
		pw, _ := cmd.Flags().GetString("password")
		if un == "" || pw == "" {
			return fmt.Errorf("--username and --password are required")
		}
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		if err := ops.SignIn(c, un, pw); err != nil {
			return err
		}
		fmt.Println("Signed in successfully")
		return nil
	},
}

var signOutCmd = &cobra.Command{
	Use:   "sign-out",
	Short: "Sign out of HEOS account",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		if err := ops.SignOut(c); err != nil {
			return err
		}
		fmt.Println("Signed out")
		return nil
	},
}

var heartBeatCmd = &cobra.Command{
	Use:   "heart-beat",
	Short: "Send heartbeat to verify connection",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		if err := ops.HeartBeat(c); err != nil {
			return err
		}
		fmt.Println("OK")
		return nil
	},
}

var rebootCmd = &cobra.Command{
	Use:   "reboot",
	Short: "Reboot the HEOS speaker",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		if err := ops.Reboot(c); err != nil {
			return err
		}
		fmt.Println("Reboot initiated")
		return nil
	},
}

var prettifyJSONCmd = &cobra.Command{
	Use:   "prettify-json",
	Short: "Enable or disable pretty JSON responses from speaker",
	RunE: func(cmd *cobra.Command, args []string) error {
		enable, _ := cmd.Flags().GetBool("enable")
		c, cleanup, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		defer cleanup()
		if err := ops.SetPrettyJSON(c, enable); err != nil {
			return err
		}
		if enable {
			fmt.Println("Pretty JSON enabled")
		} else {
			fmt.Println("Pretty JSON disabled")
		}
		return nil
	},
}
