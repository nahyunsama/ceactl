package transceiver

import (
	"crypto/tls"
	"net/http"
)

type Client struct {
	BaseURL string
	HTTP    *http.Client
	Cookie  string
}

func NewClient(host string, port string, insecureTLS bool) *Client {
	return &Client{
		BaseURL: "https://" + host + ":" + port + "/nuova",
		HTTP: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureTLS},
			},
		},
	}
}
