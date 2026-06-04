package ucsm

import "github.com/spf13/cobra"

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ucsm",
		Short: "UCSM related commands",
	}

	cmd.AddCommand(ShowServersCommand())

	return cmd
}
