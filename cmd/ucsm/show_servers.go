package ucsm

import (
	"fmt"

	"github.com/nahyunsama/ceactl/internal/ucsm/commands"
	"github.com/nahyunsama/ceactl/internal/ucsm/config"
	"github.com/spf13/cobra"
)

func ShowServersCommand(opts *commandOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "servers",
		Short: "Show UCSM Servers",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadConfig(opts.configPath, opts.deviceName, opts.verbose)
			if err != nil {
				return fmt.Errorf("failed to load config: %v", err)
			}

			servers, err := commands.GetServers(cmd.Context(), cfg)
			if err != nil {
				return fmt.Errorf("failed to get servers: %v", err)
			}

			for _, s := range servers {
				fmt.Printf("%s\t%s\t%s\t%s\n", s.DN, s.Model, s.Serial, s.OperState)
			}
			return nil
		},
	}
}
