package transceiver

import (
	"context"
	"encoding/xml"
	"fmt"
	"os"
)

type resolveClassRequest struct {
	XMLName        xml.Name `xml:"configResolveClass"`
	Cookie         string   `xml:"cookie,attr"`
	ClassID        string   `xml:"classId,attr"`
	InHierarchical string   `xml:"inHierarchical,attr"`
}

func (c *Client) ResolveClass(ctx context.Context, classID string) ([]byte, error) {
	if c.Verbose {
		fmt.Fprintf(os.Stderr, "[verbose] resolve class: %s\n", classID)
	}
	return c.PostXML(ctx, resolveClassRequest{
		Cookie:         c.Cookie,
		ClassID:        classID,
		InHierarchical: "false",
	})
}
