package transceiver

import (
	"crypto/tls"
	"net/http"
)

type Client struct {
	BaseURL  string
	HTTP     *http.Client
	Username string
	Password string
	Verbose  bool
}

func NewClient(host, port, username, password string, insecureTLS bool, verbose bool) *Client {
	return &Client{
		BaseURL:  "https://" + host + ":" + port + "/ins",
		Username: username,
		Password: password,
		HTTP: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureTLS},
			},
		},
		Verbose: verbose,
	}
}
