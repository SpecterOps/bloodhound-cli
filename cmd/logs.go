package cmd

import (
	"fmt"
	docker "github.com/SpecterOps/BloodHound_CLI/cmd/internal"
	"github.com/spf13/cobra"
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs <container>",
	Short: "Fetch logs for BloodHound services",
	Long: `Fetch logs for BloodHound services. Provide "all" or a container name.

Valid names are:

* django
* nginx
* postgres
* redis
* graphql
* queue`,
	Args: cobra.ExactArgs(1),
	Run:  readLogs,
}

func init() {
	rootCmd.AddCommand(logsCmd)

	logsCmd.Flags().StringP("lines", "l", "500", "Number of lines to display")
}

func readLogs(cmd *cobra.Command, args []string) {
	docker.EvaluateDockerComposeStatus()
	lines := cmd.Flag("lines").Value.String()
	fmt.Printf("[+] Fetching up to %s lines of logs for `%s`...\n", lines, args[0])
	logs := docker.FetchLogs(args[0], lines)
	for _, entry := range logs {
		fmt.Print(entry)
	}
}
