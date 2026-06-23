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

	client := transceiver.NewClient(cfg.SwitchIP, cfg.SwitchPort, cfg.Username, cfg.Password, cfg.InsecureTLS, cfg.Verbose)

	data, err := client.SendRequest(ctx, []byte(`{
			"ins_api": {
				"version": "1.0",
				"type": "cli_show",
				"chunk": "0",
				"sid": "1",
				"input": "show version",
				"output_format": "json"
			}
	}`))

	if err != nil {
		return receiver.Body{}, err
	}

	return receiver.ParseResponse(data)
}
