package cmd

import (
	"github.com/spf13/cobra"
)

// containersCmd represents the containers command
var containersCmd = &cobra.Command{
	Use:   "containers",
	Short: "Manage BloodHound containers with subcommands",
	Long:  "Manage BloodHound containers and services with subcommands.",
}

func init() {
	rootCmd.AddCommand(containersCmd)
}
