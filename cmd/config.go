package cmd

import (
	"fmt"
	env "github.com/SpecterOps/BloodHound_CLI/cmd/internal"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Display or adjust the configuration",
	Long: `Run this command to display the configuration. Use subcommands to
adjust the configuration or retrieve individual values.`,
	Run: configDisplay,
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func configDisplay(cmd *cobra.Command, args []string) {
	fmt.Println("[+] Current configuration and available variables:")
	configuration := env.GetConfigAll()
	fmt.Println(string(configuration))
}
