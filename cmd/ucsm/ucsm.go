package ucsm

import "github.com/spf13/cobra"

type commandOptions struct {
	configPath string
	deviceName string
	verbose    bool
}

func NewCommand() *cobra.Command {
	opts := commandOptions{}

	cmd := &cobra.Command{
		Use:   "ucsm",
		Short: "UCSM related commands",
	}

	cmd.PersistentFlags().StringVar(&opts.configPath, "config", ".config.yaml", "config file path")
	cmd.PersistentFlags().StringVarP(&opts.deviceName, "device", "d", "", "device name from config")
	cmd.PersistentFlags().BoolVarP(&opts.verbose, "verbose", "v", false, "enable verbose output")

	cmd.AddCommand(ShowServersCommand(&opts))

	return cmd
}
