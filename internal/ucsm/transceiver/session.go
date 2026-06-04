package transceiver

import (
	"context"
	"encoding/xml"
	"fmt"
)

type loginRequest struct {
	XMLName    xml.Name `xml:"aaaLogin"`
	InName     string   `xml:"inName,attr"`
	InPassword string   `xml:"inPassword,attr"`
}

type loginResponse struct {
	XMLName    xml.Name `xml:"aaaLogin"`
	OutCookie  string   `xml:"outCookie,attr"`
	ErrorCode  string   `xml:"errorCode,attr"`
	ErrorDescr string   `xml:"errorDescr,attr"`
}

func (c *Client) Login(ctx context.Context, user, password string) error {
	data, err := c.PostXML(ctx, loginRequest{InName: user, InPassword: password})
	if err != nil {
		return err
	}

	var resp loginResponse
	if err := xml.Unmarshal(data, &resp); err != nil {
		return err
	}

	if resp.ErrorCode != "" {
		return fmt.Errorf("login failed: %s", resp.ErrorDescr)
	}

	c.Cookie = resp.OutCookie
	return nil
}
