package cmd

import (
	"fmt"
	docker "github.com/SpecterOps/BloodHound_CLI/cmd/internal"
	"github.com/spf13/cobra"
)

// containersRestartCmd represents the restart command
var containersRestartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart all stopped and running BloodHound services",
	Long: `Restart all stopped and running BloodHound services. This performs
the equivalent of running the "docker compose restart" command.`,
	Run: containersRestart,
}

func init() {
	containersCmd.AddCommand(containersRestartCmd)
}

func containersRestart(cmd *cobra.Command, args []string) {
	docker.EvaluateDockerComposeStatus()
	fmt.Println("[+] Restarting the development environment")
	docker.RunDockerComposeRestart("docker-compose.yml")
}
