package cmd

import (
	"fmt"
	docker "github.com/SpecterOps/BloodHound_CLI/cmd/internal"
	"github.com/spf13/cobra"
)

// containersBuildCmd represents the build command
var containersBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Builds the BloodHound containers (only needed for updates)",
	Long: `Builds the BloodHound containers.

Note: Build will stop a container if it is already running. You will need to run
the "up" command to start the containers after the build.`,
	Run: buildContainers,
}

func init() {
	containersCmd.AddCommand(containersBuildCmd)
}

// buildContainers builds and upgrades BloodHound containers using Docker Compose.
// It checks the current Docker Compose status before initiating the build process.
func buildContainers(cmd *cobra.Command, args []string) {
	docker.EvaluateDockerComposeStatus()
	fmt.Println("[+] Starting build")
	docker.RunDockerComposeUpgrade(docker.GetYamlFilePath(dirOverride))
}
