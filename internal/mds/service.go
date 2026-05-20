// MDS Service and commands

package mds

import (
	"encoding/json"

	"github.com/nahyunsama/ceactl/internal/client"
	"github.com/nahyunsama/ceactl/internal/model"
)

type Service struct {
	client *client.NXClient
}

func NewService(c *client.NXClient) *Service {
	return &Service{client: c}
}

func (s *Service) ShowVersion() (*model.Body, error) {
	payload := `{
		"ins_api": {
			"version": "1.0",
			"type": "cli_show",
			"chunk": "0",
			"sid": "1",
			"input": "show version",
			"output_format": "json"
		}
	}`

	data, err := s.client.Post(payload)
	if err != nil {
		return nil, err
	}

	var res model.NXResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return &res.InsAPI.Outputs.Output.Body, nil
}
