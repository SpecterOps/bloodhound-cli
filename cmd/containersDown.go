package cmd

import (
	"fmt"
	docker "github.com/SpecterOps/BloodHound_CLI/cmd/internal"
	"github.com/spf13/cobra"
)

var volumes bool

// containersDownCmd represents the down command
var containersDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Bring down all BloodHound services and remove the containers",
	Long: `Bring down all BloodHound services and remove the containers. This
performs the equivalent of running the "docker compose down" command.`,
	Run: containersDown,
}

func init() {
	containersCmd.AddCommand(containersDownCmd)

	containersDownCmd.PersistentFlags().BoolVar(&volumes, "volumes", false, "Delete data volumes when containers come down")
}

func containersDown(cmd *cobra.Command, args []string) {
	docker.EvaluateDockerComposeStatus()
	fmt.Println("[+] Bringing down the BloodHound environment")
	docker.RunDockerComposeDown(docker.GetYamlFilePath(), volumes)
}
