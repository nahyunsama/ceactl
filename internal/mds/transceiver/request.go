package transceiver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type nxapiRequest struct {
	InsAPI struct {
		Version      string `json:"version"`
		Type         string `json:"type"`
		Chunk        string `json:"chunk"`
		Sid          string `json:"sid"`
		Input        string `json:"input"`
		OutputFormat string `json:"output_format"`
	} `json:"ins_api"`
}

func (c *Client) CLIShow(ctx context.Context, input string) ([]byte, error) {
	if c.Verbose {
		fmt.Fprintf(os.Stderr, "[verbose] running command: %s\n", input)
	}

	var req nxapiRequest
	req.InsAPI.Version = "1.0"
	req.InsAPI.Type = "cli_show"
	req.InsAPI.Chunk = "0"
	req.InsAPI.Sid = "1"
	req.InsAPI.Input = input
	req.InsAPI.OutputFormat = "json"

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to build nx-api request: %w", err)
	}

	return c.SendRequest(ctx, payload)
}

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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(body))
	}

	return body, nil
}
