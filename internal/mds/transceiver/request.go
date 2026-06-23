package transceiver

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/nahyunsama/ceactl/internal/mds/config"
)

func SendRequest(ctx context.Context, cfg config.Config, payload []byte) ([]byte, error) {
	url := "https://" + cfg.SwitchIP + ":" + cfg.SwitchPort + "/ins"

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(cfg.SwitchID, cfg.SwitchPW)
	req.Header.Set("Content-Type", "application/json")

	client := NewClient(cfg.InsecureTLS)

	if cfg.Verbose {
		fmt.Fprintf(os.Stderr, "[verbose] sending request: %s %s\n", req.Method, req.URL)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
