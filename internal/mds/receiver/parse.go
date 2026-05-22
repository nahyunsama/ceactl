package receiver

import "encoding/json"

func ParseResponse(data []byte) (Body, error) {
	var NXResp NXResponse
	if err := json.Unmarshal(data, &NXResp); err != nil {
		return Body{}, err
	}
	return NXResp.InsAPI.Outputs.Output.Body, nil
}
