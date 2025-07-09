package cmd

import (
	"fmt"
	docker "github.com/SpecterOps/BloodHound_CLI/cmd/internal"
	"github.com/spf13/cobra"
	"log"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Builds containers and performs first-time setup of BloodHound",
	Long: `Builds containers and performs first-time setup of BloodHound.

The command performs the following steps:

* Sets up the default server configuration
* Builds the Docker containers
* Creates a default admin user with a randomly generated password

This command only needs to be run once. If you run it again, you will see some errors because
certain actions (e.g., creating the default user) can and should only be done once.`,
	Run: installBloodHound,
}

// init registers the install command with the root command, making it available in the CLI.
func init() {
	rootCmd.AddCommand(installCmd)
}

// installBloodHound sets up the BloodHound environment by verifying Docker Compose status, creating the required home directory, and launching the Docker containers using the installation configuration.
func installBloodHound(cmd *cobra.Command, args []string) {
	docker.EvaluateDockerComposeStatus()
	configErr := docker.MakeConfigDir()
	if configErr != nil {
		log.Fatalf("Error creating config directory: %v", configErr)
	}
	fmt.Println("[+] Starting BloodHound environment installation")
	docker.RunDockerComposeInstall(docker.GetYamlFilePath(dirOverride))
}
