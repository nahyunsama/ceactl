package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/nahyunsama/ceactl/internal/mds/config"
	"github.com/nahyunsama/ceactl/internal/mds/receiver"
	"github.com/nahyunsama/ceactl/internal/mds/transceiver"
)

func GetInventory(ctx context.Context, cfg config.Config) (receiver.InventoryBody, error) {
	if cfg.Verbose {
		fmt.Fprintf(os.Stderr, "[verbose] running command: show inventory\n")
	}

	client := transceiver.NewClient(cfg.SwitchIP, cfg.SwitchPort, cfg.Username, cfg.Password, cfg.InsecureTLS, cfg.Verbose)

	data, err := client.SendRequest(ctx, []byte(`{
			"ins_api": {
				"version": "1.0",
				"type": "cli_show",
				"chunk": "0",
				"sid": "1",
				"input": "show inventory",
				"output_format": "json"
			}
	}`))

	if err != nil {
		return receiver.InventoryBody{}, err
	}

	return receiver.ParseInventoryResponse(data)
}
