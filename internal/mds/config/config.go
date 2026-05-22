package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	SwitchID    string
	SwitchPW    string
	SwitchIP    string
	SwitchPort  string
	InsecureTLS bool
}

func LoadConfig() (Config, error) {
	if err := godotenv.Load(); err != nil {
		return Config{}, fmt.Errorf("error loading .env file: %w", err)
	}

	cfg := Config{
		SwitchID:    os.Getenv("switch_ID"),
		SwitchPW:    os.Getenv("switch_PW"),
		SwitchIP:    os.Getenv("switch_IP"),
		SwitchPort:  os.Getenv("switch_Port"),
		InsecureTLS: true,
	}

	if cfg.SwitchID == "" || cfg.SwitchPW == "" || cfg.SwitchIP == "" || cfg.SwitchPort == "" {
		return Config{}, fmt.Errorf("missing required environment variables")
	}

	return cfg, nil
}
