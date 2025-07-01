package cmd

import (
	"fmt"
	docker "github.com/SpecterOps/BloodHound_CLI/cmd/internal"
	"github.com/spf13/cobra"
	"log"
	"path/filepath"
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

func init() {
	rootCmd.AddCommand(installCmd)
}

func installBloodHound(cmd *cobra.Command, args []string) {
	err := docker.EvaluateDockerComposeStatus()
	if err != nil {
		return
	}
	homeErr := docker.MakeHomeDir()
	if homeErr != nil {
		log.Fatalf("Error creating home directory: %v", homeErr)
	}
	fmt.Println("[+] Starting BloodHound environment installation")
	docker.RunDockerComposeInstall(filepath.Join(docker.GetBloodHoundDir(), "docker-compose.yml"))
}
