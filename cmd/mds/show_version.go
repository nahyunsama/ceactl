package mds

import (
	"context"
	"fmt"
	"log"

	"github.com/nahyunsama/ceactl/internal/mds/commands"
	"github.com/nahyunsama/ceactl/internal/mds/config"
	"github.com/spf13/cobra"
)

func ShowVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show MDS Firmware Version",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.LoadConfig()
			if err != nil {
				log.Fatalf("Failed to load config: %v", err)
			}

			info, err := commands.GetVersion(context.Background(), cfg)
			if err != nil {
				log.Fatalf("Failed to get version: %v", err)
			}

			fmt.Printf("Host Name: %s\n", info.HostName)
			fmt.Printf("Version: %s\n", info.Version)
		},
	}
}
