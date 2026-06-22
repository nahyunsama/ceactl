package config

import appconfig "github.com/nahyunsama/ceactl/internal/config"

type Config struct {
	UCSMID      string
	UCSMPW      string
	UCSMIP      string
	UCSMPort    string
	InsecureTLS bool
}

func LoadConfig(configPath, deviceName string) (Config, error) {
	device, err := appconfig.LoadDevice(configPath, deviceName, "ucsm")
	if err != nil {
		return Config{}, err
	}

	return Config{
		UCSMID:      device.Username,
		UCSMPW:      device.Password,
		UCSMIP:      device.Host,
		UCSMPort:    device.Port,
		InsecureTLS: device.InsecureTLS,
	}, nil
}
