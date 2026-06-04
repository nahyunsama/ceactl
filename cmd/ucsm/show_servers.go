package ucsm

import (
	"context"
	"fmt"
	"log"

	"github.com/nahyunsama/ceactl/internal/ucsm/commands"
	"github.com/nahyunsama/ceactl/internal/ucsm/config"
	"github.com/spf13/cobra"
)

func ShowServersCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "servers",
		Short: "Show UCSM Servers",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.LoadConfig()
			if err != nil {
				log.Fatalf("Failed to load config: %v", err)
			}

			servers, err := commands.GetServers(context.Background(), cfg)
			if err != nil {
				log.Fatalf("Failed to get servers: %v", err)
			}

			for _, s := range servers {
				fmt.Printf("%s\t%s\t%s\t%s\n", s.DN, s.Model, s.Serial, s.OperState)
			}
		},
	}
}
