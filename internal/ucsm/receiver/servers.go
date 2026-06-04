package receiver

import "encoding/xml"

type ServersResponse struct {
	XMLName    xml.Name `xml:"configResolveClass"`
	OutConfigs struct {
		Blades    []Server `xml:"computeBlade"`
		RackUnits []Server `xml:"computeRackUnit"`
	} `xml:"outConfigs"`
}

type Server struct {
	DN        string `xml:"dn,attr"`
	Model     string `xml:"model,attr"`
	Serial    string `xml:"serial,attr"`
	OperState string `xml:"operState,attr"`
}

func ParseServers(data []byte) ([]Server, error) {
	var resp ServersResponse
	if err := xml.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	servers := append(resp.OutConfigs.RackUnits, resp.OutConfigs.Blades...)
	return servers, nil
}
