package cmd

import (
	"fmt"
	docker "github.com/SpecterOps/BloodHound_CLI/cmd/internal"
	"github.com/spf13/cobra"
	"path/filepath"
)

// installCmd represents the install command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove all BloodHound containers, images, and volume data",
	Long: `Remove all BloodHound containers, images, and volume data.

The command performs the following steps:

* Brings down running containers
* Deletes the stopped containers
* Deletes the container images
* Deletes all BloodHound volumes and data

This command is irreversible and should only be run if you are looking to remove BloodHound from the system or wanting
a fresh start.`,
	Run: uninstallBloodHound,
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}

// uninstallBloodHound removes all BloodHound Docker containers, images, and volumes when the uninstall command is invoked.
// It first checks the Docker Compose environment status and proceeds with the uninstallation if no errors are detected.
func uninstallBloodHound(cmd *cobra.Command, args []string) {
	err := docker.EvaluateDockerComposeStatus()
	if err != nil {
		return
	}
	fmt.Println("[+] Starting BloodHound environment removal")
	docker.RunDockerComposeUninstall(filepath.Join(docker.GetBloodHoundDir(), "docker-compose.yml"))
}
