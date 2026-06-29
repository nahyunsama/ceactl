package mds

import (
	"fmt"

	"github.com/nahyunsama/ceactl/internal/mds/commands"
	"github.com/nahyunsama/ceactl/internal/mds/config"
	"github.com/spf13/cobra"
)

func ShowVersionCommand(opts *commandOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show MDS Firmware Version",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadConfig(opts.configPath, opts.deviceName, opts.verbose)
			if err != nil {
				return fmt.Errorf("failed to load config: %v", err)
			}

			info, err := commands.GetVersion(cmd.Context(), cfg)
			if err != nil {
				return fmt.Errorf("failed to get version: %v", err)
			}

			fmt.Printf("Host Name: %s\n", info.HostName)
			fmt.Printf("Version: %s\n", info.Version)
			fmt.Printf("Uptime : %d days %02d:%02d:%02d\n", info.UptimeDays, info.UptimeHrs, info.UptimeMins, info.UptimeSecs)
			return nil
		},
	}
}
