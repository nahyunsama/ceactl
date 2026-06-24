package receiver

import "encoding/json"

func ParseVersionResponse(data []byte) (VersionBody, error) {
	var resp VersionResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return VersionBody{}, err
	}
	return resp.InsAPI.Outputs.Output.Body, nil
}

func ParseInventoryResponse(data []byte) (InventoryBody, error) {
	var resp InventoryResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return InventoryBody{}, err
	}
	return resp.InsAPI.Outputs.Output.Body, nil
}
