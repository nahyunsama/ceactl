package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	UCSMID      string
	UCSMPW      string
	UCSMIP      string
	UCSMPort    string
	InsecureTLS bool
}

func LoadConfig() (Config, error) {
	if err := godotenv.Load(); err != nil {
		return Config{}, fmt.Errorf("error loading .env file: %w", err)
	}

	cfg := Config{
		UCSMID:      os.Getenv("UCSM_ID"),
		UCSMPW:      os.Getenv("UCSM_PW"),
		UCSMIP:      os.Getenv("UCSM_IP"),
		UCSMPort:    os.Getenv("UCSM_Port"),
		InsecureTLS: true,
	}
	if cfg.UCSMID == "" || cfg.UCSMPW == "" || cfg.UCSMIP == "" || cfg.UCSMPort == "" {
		return Config{}, fmt.Errorf("missing required environment variables")
	}

	return cfg, nil
}
