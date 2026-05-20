package client

import (
	"bytes"
	"crypto/tls"
	"io"
	"net/http"
)

const insecureTLS = true

type NXClient struct {
	BaseURL  string
	Username string
	Password string
	Client   *http.Client
}

func New(baseURL, username, password string) *NXClient {
	return &NXClient{
		BaseURL:  baseURL,
		Username: username,
		Password: password,
		Client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: insecureTLS,
				},
			},
		},
	}
}

func (c *NXClient) Post(payload string) ([]byte, error) {
	req, err := http.NewRequest(
		"POST",
		c.BaseURL+"/ins",
		bytes.NewBuffer([]byte(payload)),
	)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.Username, c.Password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
