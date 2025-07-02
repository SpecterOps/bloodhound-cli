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

// init registers the version command with the root command, enabling the "version" CLI command.
func init() {
	rootCmd.AddCommand(versionCmd)
}

// displayVersion prints the BloodHound CLI version and build date to standard output.
func displayVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("BloodHound-CLI (%s, %s)\n", config.Version, config.BuildDate)
}
