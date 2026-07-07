package commands

import (
	"context"

	"github.com/nahyunsama/ceactl/internal/mds/config"
	"github.com/nahyunsama/ceactl/internal/mds/receiver"
	"github.com/nahyunsama/ceactl/internal/mds/transceiver"
)

func GetInventory(ctx context.Context, cfg config.Config) (receiver.InventoryBody, error) {
	client := transceiver.NewClientFromConfig(cfg)

	data, err := client.CLIShow(ctx, "show inventory")
	if err != nil {
		return receiver.InventoryBody{}, err
	}

	return receiver.ParseInventoryResponse(data)
}
