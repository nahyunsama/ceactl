package mds

import "github.com/spf13/cobra"

type commandOptions struct {
	configPath string
	deviceName string
}

func NewCommand() *cobra.Command {
	opts := commandOptions{}

	cmd := &cobra.Command{
		Use:   "mds",
		Short: "MDS related commands",
	}

	cmd.PersistentFlags().StringVar(&opts.configPath, "config", ".config.yaml", "config file path")
	cmd.PersistentFlags().StringVarP(&opts.deviceName, "device", "d", "", "device name from config")

	cmd.AddCommand(ShowVersionCommand(&opts))

	return cmd
}
