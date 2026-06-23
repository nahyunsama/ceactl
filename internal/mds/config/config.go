package config

import appconfig "github.com/nahyunsama/ceactl/internal/config"

type Config struct {
	SwitchIP    string
	SwitchPort  string
	Username    string
	Password    string
	InsecureTLS bool
	Verbose     bool
}

func LoadConfig(configPath, deviceName string, verbose bool) (Config, error) {
	device, err := appconfig.LoadDevice(configPath, deviceName, "mds")
	if err != nil {
		return Config{}, err
	}

	return Config{
		SwitchIP:    device.Host,
		SwitchPort:  device.Port,
		Username:    device.Username,
		Password:    device.Password,
		InsecureTLS: device.InsecureTLS,
		Verbose:     verbose,
	}, nil
}
