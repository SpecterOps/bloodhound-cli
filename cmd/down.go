package cmd

import (
	"github.com/spf13/cobra"
)

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Shortcut for `containers down`",
	Run: func(cmd *cobra.Command, args []string) {
		containersDownCmd.Run(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
}
