package transceiver

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
)

func (c *Client) SendRequest(ctx context.Context, payload []byte) ([]byte, error) {

	req, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.Username, c.Password)

	if c.Verbose {
		fmt.Fprintf(os.Stderr, "[verbose] sending request: %s %s\n", req.Method, c.BaseURL)
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
