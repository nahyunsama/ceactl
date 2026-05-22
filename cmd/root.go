package cmd

import (
	"os"

	"github.com/nahyunsama/ceactl/cmd/mds"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ceactl",
	Short: "Cisco Enterprise API Control CLI",
}

func init() {
	rootCmd.AddCommand(mds.NewCommand())
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
