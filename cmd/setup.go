package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jrogala/heos-cli/client"
	"github.com/jrogala/heos-cli/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(setupCmd)
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Interactive setup to configure HEOS speaker connection",
	RunE: func(cmd *cobra.Command, args []string) error {
		reader := bufio.NewReader(os.Stdin)

		cfg := config.Get()
		defaultHost := cfg.Host
		defaultPort := cfg.Port
		if defaultPort == 0 {
			defaultPort = config.DefaultPort
		}

		fmt.Printf("HEOS speaker IP address [%s]: ", defaultHost)
		host, _ := reader.ReadString('\n')
		host = strings.TrimSpace(host)
		if host == "" {
			host = defaultHost
		}
		if host == "" {
			return fmt.Errorf("host is required")
		}

		fmt.Printf("HEOS port [%d]: ", defaultPort)
		portStr, _ := reader.ReadString('\n')
		portStr = strings.TrimSpace(portStr)
		port := defaultPort
		if portStr != "" {
			p, err := strconv.Atoi(portStr)
			if err != nil {
				return fmt.Errorf("invalid port: %s", portStr)
			}
			port = p
		}

		fmt.Printf("Testing connection to %s:%d... ", host, port)
		c := client.New(host, port)
		if err := c.Connect(); err != nil {
			fmt.Println("FAILED")
			return err
		}
		if err := c.HeartBeat(); err != nil {
			c.Close()
			fmt.Println("FAILED")
			return fmt.Errorf("heartbeat failed: %w", err)
		}
		c.Close()
		fmt.Println("OK")

		if err := config.Save(host, port); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}
		fmt.Printf("Configuration saved to %s/config.yaml\n", config.ConfigDir())
		return nil
	},
}
