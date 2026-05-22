// MDS Login and print inventory
// error handing, logging, and configuration management
//

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/nahyunsama/ceactl/internal/mds/commands"
	"github.com/nahyunsama/ceactl/internal/mds/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	info, err := commands.GetVersion(context.Background(), cfg)
	if err != nil {
		log.Fatalf("Failed to get version: %v", err)
	}

	fmt.Printf("Host Name: %s\n", info.HostName)
	fmt.Printf("Version: %s\n", info.Version)
}
