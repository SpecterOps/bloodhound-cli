package cmd

import (
	"github.com/spf13/cobra"
)

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Shortcut for `containers up`",
	Run: func(cmd *cobra.Command, args []string) {
		containersUpCmd.Run(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}
