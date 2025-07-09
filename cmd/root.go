package cmd

import (
	env "github.com/SpecterOps/BloodHound_CLI/cmd/internal"
	"github.com/spf13/cobra"
	"os"
)

// Vars for global flags
var fileOverride string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bloodhound-cli",
	Short: "A command line interface for managing BloodHound.",
	Long: `BloodHound CLI is a command line interface for managing BloodHound and
associated containers and services. Commands are grouped by their use.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Create or parse the Docker ``bloodhound.config.json`` file
	env.ParseBloodHoundEnvironmentVariables()

	rootCmd.PersistentFlags().StringVarP(&fileOverride, "file", "f", "", `Override the YAML file in the configured data directory and use a different YAML file for the container commands.`)
}
