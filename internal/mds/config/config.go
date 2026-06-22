package config

import appconfig "github.com/nahyunsama/ceactl/internal/config"

type Config struct {
	SwitchID    string
	SwitchPW    string
	SwitchIP    string
	SwitchPort  string
	InsecureTLS bool
}

func LoadConfig(configPath, deviceName string) (Config, error) {
	device, err := appconfig.LoadDevice(configPath, deviceName, "mds")
	if err != nil {
		return Config{}, err
	}

	return Config{
		SwitchID:    device.Username,
		SwitchPW:    device.Password,
		SwitchIP:    device.Host,
		SwitchPort:  device.Port,
		InsecureTLS: device.InsecureTLS,
	}, nil
}
