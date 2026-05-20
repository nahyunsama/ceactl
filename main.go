// MDS Login and print inventory

package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

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

	// fmt.Printf("Switch ID: %s\n", switchID)
	// fmt.Printf("Switch Password: %s\n", switchPW)
	// fmt.Printf("Switch IP: %s\n", switchIP)
	// fmt.Printf("Switch Port: %s\n", switchPort)

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

	fmt.Printf("Response: %s\n", string(body))
}
