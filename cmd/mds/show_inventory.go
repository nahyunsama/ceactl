package mds

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/nahyunsama/ceactl/internal/mds/commands"
	"github.com/nahyunsama/ceactl/internal/mds/config"
	"github.com/spf13/cobra"
)

func ShowInventoryCommand(opts *commandOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "inventory",
		Short: "Show MDS Inventory",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadConfig(opts.configPath, opts.deviceName, opts.verbose)
			if err != nil {
				return fmt.Errorf("failed to load config: %v", err)
			}

			info, err := commands.GetInventory(context.Background(), cfg)
			if err != nil {
				return fmt.Errorf("failed to get inventory: %v", err)
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintf(w, "NAME\tPRODUCT ID\tSERIAL NUM\n")
			for _, item := range info.TableInv.RowInv {
				fmt.Fprintf(w, "%s\t%s\t%s\n", item.Name, item.ProductId, item.SerialNum)
			}
			w.Flush()

			return nil
		},
	}
}
