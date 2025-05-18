package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	docker "github.com/SpecterOps/BloodHound_CLI/cmd/internal"
	"github.com/spf13/cobra"
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

func uninstallBloodHound(cmd *cobra.Command, args []string) {
	err := docker.EvaluateDockerComposeStatus(true)
	if err != nil {
		return
	}
	fmt.Println("[+] Starting BloodHound environment removal")

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("[-] Failed to get current directory:", err)
		return
	}

	// Build full path to docker-compose.yml
	composePath := filepath.Join(cwd, "docker-compose.yml")

	docker.RunDockerComposeUninstall(composePath)
}
