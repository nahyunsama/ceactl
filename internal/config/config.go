package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	SwitchID   string
	SwitchPW   string
	SwitchIP   string
	SwitchPort string
}

func Load() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env not found, using system env")
	}

	return Config{
		SwitchID:   os.Getenv("switch_ID"),
		SwitchPW:   os.Getenv("switch_PW"),
		SwitchIP:   os.Getenv("switch_IP"),
		SwitchPort: os.Getenv("switch_Port"),
	}
}
