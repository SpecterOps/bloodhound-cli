package cmd

import (
	"fmt"
	docker "github.com/SpecterOps/BloodHound_CLI/cmd/internal"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Evaluates the Docker environment and downloads the necessary YAML files, as needed.",
	Long: `Evaluates the Docker environment and downloads the necessary YAML files, as needed.

You can run this command before or after running the "install" command. The intent is to ensure that
the necessary commands are available in the $PATH and the YAML files are downloaded. If you accidentally delete the
YAML files or move the binary without them, this command will prompt you to re-download them.`,
	Run: evaluateBloodHound,
}

func init() {
	rootCmd.AddCommand(checkCmd)
}

func evaluateBloodHound(cmd *cobra.Command, args []string) {
	err := docker.EvaluateDockerComposeStatus()
	if err != nil {
		return
	}
	docker.EvaluateEnvironment()
	fmt.Println("[+] Environment checks are complete!")
}
