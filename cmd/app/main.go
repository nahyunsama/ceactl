// MDS Login and print inventory
// error handing, logging, and configuration management
//

package main

import (
	"fmt"
	"log"
	"net/url"

	"github.com/nahyunsama/ceactl/internal/client"
	"github.com/nahyunsama/ceactl/internal/config"
	"github.com/nahyunsama/ceactl/internal/mds"
)

func main() {

	cfg := config.Load()

	baseURL := fmt.Sprintf(
		"https://%s:%s",
		cfg.SwitchIP,
		cfg.SwitchPort,
	)

	_, err := url.Parse(baseURL)
	if err != nil {
		log.Fatal(err)
	}

	c := client.New(baseURL, cfg.SwitchID, cfg.SwitchPW)
	svc := mds.NewService(c)

	res, err := svc.ShowVersion()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("HostName: %s\n", res.HostName)
	fmt.Printf("Version: %s\n", res.Version)
}
