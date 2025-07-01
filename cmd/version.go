package cmd

import (
	"fmt"
	"github.com/SpecterOps/BloodHound_CLI/cmd/config"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Displays BloodHound CLI's version information",
	Long:  "Displays BloodHound CLI's version information.",
	Run:   displayVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func displayVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("BloodHound-CLI (%s, %s)\n", config.Version, config.BuildDate)
}
