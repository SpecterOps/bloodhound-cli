package cmd

import (
	"fmt"
	docker "github.com/SpecterOps/BloodHound_CLI/cmd/internal"
	"github.com/spf13/cobra"
)

// containersStopCmd represents the stop command
var containersStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop all BloodHound services without removing the containers",
	Long: `Stop all BloodHound services without removing the containers. This
performs the equivalent of running the "docker compose stop" command.`,
	Run: containersStop,
}

func init() {
	containersCmd.AddCommand(containersStopCmd)
}

// containersStop stops all BloodHound Docker Compose services without removing their containers.
func containersStop(cmd *cobra.Command, args []string) {
	docker.EvaluateDockerComposeStatus()
	fmt.Println("[+] Stopping the BloodHound environment")
	docker.RunDockerComposeStop(docker.GetYamlFilePath(fileOverride))
}
