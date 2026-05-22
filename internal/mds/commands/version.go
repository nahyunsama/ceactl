package commands

import (
	"context"

	"github.com/nahyunsama/ceactl/internal/mds/config"
	"github.com/nahyunsama/ceactl/internal/mds/receiver"
	"github.com/nahyunsama/ceactl/internal/mds/transceiver"
)

func GetVersion(ctx context.Context, cfg config.Config) (receiver.Body, error) {
	payload := []byte(`{
		"ins_api": {
			"version": "1.0",
			"type": "cli_show",
			"chunk": "0",
			"sid": "1",
			"input": "show version",
			"output_format": "json"
		}
	}`)

	data, err := transceiver.SendRequest(ctx, cfg, payload)
	if err != nil {
		return receiver.Body{}, err
	}

	return receiver.ParseResponse(data)
}
