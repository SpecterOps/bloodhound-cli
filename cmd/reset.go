package cmd

import (
	"fmt"
	docker "github.com/SpecterOps/BloodHound_CLI/cmd/internal"
	"github.com/spf13/cobra"
)

// resetPwdCmd represents the resetpwd command
var resetPwdCmd = &cobra.Command{
	Use:   "resetpwd",
	Short: "Reset the admin password",
	Long: `Reset the admin password by bringing down any running containers and recreating the default admin record.

The command performs the following steps:

* Brings down any running containers
* Temporarily sets the "bhe_recreate_default_admin" environment variable to "true"
* Generates a new password
* Brings the containers back up to recreate the admin record

**WARNING** : This action wipes all user data for the default admin user. This action cannot be undone.
`,
	Run: resetAdminPwd,
}

func init() {
	rootCmd.AddCommand(resetPwdCmd)
}

func resetAdminPwd(cmd *cobra.Command, args []string) {
	err := docker.EvaluateDockerComposeStatus(true)
	if err != nil {
		return
	}
	fmt.Println("[+] Resetting admin password")
	docker.ResetAdminPassword("docker-compose.yml")
}
