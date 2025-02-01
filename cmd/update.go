package cmd

import (
	"fmt"
	docker "github.com/SpecterOps/BloodHound_CLI/cmd/internal"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the BloodHound container images if an update is available",
	Long:  `Updates the BloodHound container images if an update is available.`,
	Run:   updateBloodHound,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func updateBloodHound(cmd *cobra.Command, args []string) {
	err := docker.EvaluateDockerComposeStatus()
	if err != nil {
		return
	}
	fmt.Println("[+] Checking for BloodHound image updates...")
	docker.RunDockerComposePull("docker-compose.yml")
}
