package transceiver

import (
	"context"
	"encoding/xml"
)

type resolveClassRequest struct {
	XMLName        xml.Name `xml:"configResolveClass"`
	Cookie         string   `xml:"cookie,attr"`
	ClassID        string   `xml:"classId,attr"`
	InHierarchical string   `xml:"inHierarchical,attr"`
}

func (c *Client) ResolveClass(ctx context.Context, classID string) ([]byte, error) {
	return c.PostXML(ctx, resolveClassRequest{
		Cookie:         c.Cookie,
		ClassID:        classID,
		InHierarchical: "false",
	})
}
