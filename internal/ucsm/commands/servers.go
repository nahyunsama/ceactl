package commands

import (
	"context"

	"github.com/nahyunsama/ceactl/internal/ucsm/config"
	"github.com/nahyunsama/ceactl/internal/ucsm/receiver"
	"github.com/nahyunsama/ceactl/internal/ucsm/transceiver"
)

func GetServers(ctx context.Context, cfg config.Config) ([]receiver.Server, error) {
	client := transceiver.NewClient(cfg.UCSMIP, cfg.UCSMPort, cfg.InsecureTLS)

	if err := client.Login(ctx, cfg.UCSMID, cfg.UCSMPW); err != nil {
		return nil, err
	}

	data, err := client.ResolveClass(ctx, "computeItem")
	if err != nil {
		return nil, err
	}

	return receiver.ParseServers(data)
}
