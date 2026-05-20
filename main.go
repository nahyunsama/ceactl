// MDS Login and print inventory
// error handing, logging, and configuration management
//

package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type NXResponse struct {
	InsAPI struct {
		Outputs struct {
			Output struct {
				Body Body `json:"body"`
			} `json:"output"`
		} `json:"outputs"`
	} `json:"ins_api"`
}

type Body struct {
	HostName string `json:"host_name"`
	Version  string `json:"sys_ver_str"`
}

const insecureTLS = true

func main() {
	fmt.Println("Hello CeaCtl")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	switchID := os.Getenv("switch_ID")
	switchPW := os.Getenv("switch_PW")
	switchIP := os.Getenv("switch_IP")
	switchPort := os.Getenv("switch_Port")

	url := "https://" + switchIP + ":" + switchPort + "/ins"

	payload := `{
		"ins_api": {
			"version": "1.0",
			"type": "cli_show",
			"chunk": "0",
			"sid": "1",
			"input": "show version",
			"output_format": "json"
		}

	}`

	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer([]byte(payload)),
	)
	if err != nil {
		panic(err)
	}

	req.SetBasicAuth(switchID, switchPW)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: insecureTLS,
			},
		},
	}

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	var nxResponse NXResponse
	if err := json.Unmarshal(body, &nxResponse); err != nil {
		log.Fatal("Error unmarshaling JSON response")
	}

	switchInfo := nxResponse.InsAPI.Outputs.Output.Body

	fmt.Printf("Host Name: %s\n", switchInfo.HostName)
	fmt.Printf("Version: %s\n", switchInfo.Version)
}
