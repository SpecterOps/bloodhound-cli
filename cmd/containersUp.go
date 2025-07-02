package cmd

import (
	"fmt"
	docker "github.com/SpecterOps/BloodHound_CLI/cmd/internal"
	"github.com/spf13/cobra"
)

// containersUpCmd represents the up command
var containersUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Build, (re)create, and start all BloodHound containers",
	Long: `Build, (re)create, and start all BloodHound containers. This
performs the equivalent of running the "docker compose up" command.`,
	Run: containersUp,
}

func init() {
	containersCmd.AddCommand(containersUpCmd)
}

func containersUp(cmd *cobra.Command, args []string) {
	docker.EvaluateDockerComposeStatus()
	fmt.Println("[+] Bringing up the BloodHound environment")
	docker.RunDockerComposeUp(docker.GetYamlFilePath())
}
