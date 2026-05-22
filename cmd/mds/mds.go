package mds

import "github.com/spf13/cobra"

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mds",
		Short: "MDS related commands",
	}

	cmd.AddCommand(ShowVersionCommand())

	return cmd
}
