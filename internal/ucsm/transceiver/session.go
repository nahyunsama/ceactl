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

type logoutRequest struct {
	XMLName  xml.Name `xml:"aaaLogout"`
	InCookie string   `xml:"inCookie,attr"`
}

type logoutResponse struct {
	XMLName    xml.Name `xml:"aaaLogout"`
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

func (c *Client) Logout(ctx context.Context) error {
	if c.Cookie == "" {
		return nil
	}

	data, err := c.PostXML(ctx, logoutRequest{InCookie: c.Cookie})
	if err != nil {
		return err
	}

	var resp logoutResponse
	if err := xml.Unmarshal(data, &resp); err != nil {
		return err
	}

	if resp.ErrorCode != "" {
		return fmt.Errorf("logout failed: %s", resp.ErrorDescr)
	}

	c.Cookie = ""
	return nil
}
