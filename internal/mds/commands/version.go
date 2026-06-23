package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/nahyunsama/ceactl/internal/mds/config"
	"github.com/nahyunsama/ceactl/internal/mds/receiver"
	"github.com/nahyunsama/ceactl/internal/mds/transceiver"
)

func GetVersion(ctx context.Context, cfg config.Config) (receiver.Body, error) {
	if cfg.Verbose {
		fmt.Fprintf(os.Stderr, "[verbose] running command: show version\n")
	}
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

	if cfg.Verbose {
		fmt.Fprintf(os.Stderr, "[verbose] sending request: %s\n", string(payload))
	}

	data, err := transceiver.SendRequest(ctx, cfg, payload)
	if err != nil {
		return receiver.Body{}, err
	}

	return receiver.ParseResponse(data)
}
