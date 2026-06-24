package cmd

import (
	"os"

	"github.com/nahyunsama/ceactl/cmd/mds"
	"github.com/nahyunsama/ceactl/cmd/ucsm"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "ceactl",
	Short:        "Cisco Enterprise API Control CLI",
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(mds.NewCommand())
	rootCmd.AddCommand(ucsm.NewCommand())
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
