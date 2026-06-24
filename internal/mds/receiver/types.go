package receiver

type VersionResponse struct {
	InsAPI struct {
		Outputs struct {
			Output struct {
				Body VersionBody `json:"body"`
			} `json:"output"`
		} `json:"outputs"`
	} `json:"ins_api"`
}

type VersionBody struct {
	HostName string `json:"host_name"`
	Version  string `json:"sys_ver_str"`
}

type InventoryResponse struct {
	InsAPI struct {
		Outputs struct {
			Output struct {
				Body InventoryBody `json:"body"`
			} `json:"output"`
		} `json:"outputs"`
	} `json:"ins_api"`
}

type InventoryItem struct {
	Name      string `json:"name"`
	ProductId string `json:"productid"`
	SerialNum string `json:"serialnum"`
}

type InventoryBody struct {
	TableInv struct {
		RowInv []InventoryItem `json:"ROW_inv"`
	} `json:"TABLE_inv"`
}
