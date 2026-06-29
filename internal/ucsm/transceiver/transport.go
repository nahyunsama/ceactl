package transceiver

import (
	"crypto/tls"
	"net/http"
	"time"
)

type Client struct {
	BaseURL string
	HTTP    *http.Client
	Cookie  string
	Verbose bool
}

func NewClient(host string, port string, insecureTLS bool, verbose bool) *Client {
	return &Client{
		BaseURL: "https://" + host + ":" + port + "/nuova",
		HTTP: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureTLS},
			},
		},
		Verbose: verbose,
	}
}
