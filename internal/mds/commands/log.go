package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/nahyunsama/ceactl/internal/mds/config"
	"github.com/nahyunsama/ceactl/internal/mds/receiver"
	"github.com/nahyunsama/ceactl/internal/mds/transceiver"
)

func GetLoggingLogfile(ctx context.Context, cfg config.Config) (string, error) {
	if cfg.Verbose {
		fmt.Fprintf(os.Stderr, "[verbose] running command: show logging logfile\n")
	}

	client := transceiver.NewClient(cfg.SwitchIP, cfg.SwitchPort, cfg.Username, cfg.Password, cfg.InsecureTLS, cfg.Verbose)
	data, err := client.SendRequest(ctx, []byte(`{
			"ins_api": {
				"version": "1.0",
				"type": "cli_show",
				"chunk": "0",
				"sid": "1",
				"input": "show logging logfile",
				"output_format": "json"
			}
	}`))

	if err != nil {
		return "", err
	}
	return receiver.ParseLoggingResponse(data)
}
