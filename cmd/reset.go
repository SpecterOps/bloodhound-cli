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

**NOTE** : This command requires BloodHound >= v7.1.0.

**WARNING** : This action wipes all user data for the default admin user. This action cannot be undone.
`,
	Run: resetAdminPwd,
}

// init registers the resetpwd command with the root command for the CLI.
func init() {
	rootCmd.AddCommand(resetPwdCmd)
}

// resetAdminPwd resets the default admin password by orchestrating Docker Compose operations and invoking the password reset process.
func resetAdminPwd(cmd *cobra.Command, args []string) {
	docker.EvaluateDockerComposeStatus()
	fmt.Println("[+] Resetting admin password")
	docker.ResetAdminPassword(docker.GetYamlFilePath())
}
