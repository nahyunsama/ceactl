package model

type NXResponse struct {
	InsAPI struct {
		Outputs struct {
			Output struct {
				Body Body `json:"body"`
			} `json:"output"`
		} `json:"outputs"`
	} `json:"ins_api"`
}

type Body struct {
	HostName string `json:"host_name"`
	Version  string `json:"sys_ver_str"`
}
