package transceiver

import (
	"crypto/tls"
	"net/http"
)

func NewClient(insecureTLS bool) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: insecureTLS,
			},
		},
	}
}
