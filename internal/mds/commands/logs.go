package commands

import (
	"context"

	"github.com/nahyunsama/ceactl/internal/mds/config"
	"github.com/nahyunsama/ceactl/internal/mds/receiver"
	"github.com/nahyunsama/ceactl/internal/mds/transceiver"
)

func GetLoggingLogfile(ctx context.Context, cfg config.Config) (string, error) {
	client := transceiver.NewClientFromConfig(cfg)

	data, err := client.CLIShow(ctx, "show logging logfile")
	if err != nil {
		return "", err
	}
	return receiver.ParseLoggingResponse(data)
}
