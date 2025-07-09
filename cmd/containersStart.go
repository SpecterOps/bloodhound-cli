package cmd

import (
	"fmt"
	docker "github.com/SpecterOps/BloodHound_CLI/cmd/internal"
	"github.com/spf13/cobra"
)

// containersStartCmd represents the start command
var containersStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start all stopped BloodHound services",
	Long: `Start all stopped BloodHound services. This performs the equivalent
of running the "docker compose start" command.`,
	Run: containersStart,
}

func init() {
	containersCmd.AddCommand(containersStartCmd)
}

// containersStart starts all stopped BloodHound Docker Compose services by invoking the appropriate Docker Compose command with the configuration file located in the BloodHound directory.
func containersStart(cmd *cobra.Command, args []string) {
	docker.EvaluateDockerComposeStatus()
	fmt.Println("[+] Starting the BloodHound environment")
	docker.RunDockerComposeStart(docker.GetYamlFilePath(dirOverride))
}
