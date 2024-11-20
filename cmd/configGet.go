package cmd

import (
	"fmt"
	env "github.com/SpecterOps/BloodHound_CLI/cmd/internal"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"text/tabwriter"
)

// configGetCmd represents the configGet command
var configGetCmd = &cobra.Command{
	Use:   "get <configuration> <configuration> ...",
	Short: "Get the specified configuration values",
	Long: `Get the specified configuration values. You can provide one value or
a list of values separated by spaces.

For example: bloodhound-cli config get ADMIN_PASSWORD POSTGRES_PASSWORD`,
	Run: configGet,
}

func init() {
	configCmd.AddCommand(configGetCmd)
}

func configGet(cmd *cobra.Command, args []string) {
	// initialize tabwriter
	writer := new(tabwriter.Writer)
	// Set minwidth, tabwidth, padding, padchar, and flags
	writer.Init(os.Stdout, 8, 8, 1, '\t', 0)

	defer writer.Flush()

	fmt.Println("[+] Getting configuration values:")
	fmt.Fprintf(writer, "\n %s\t%s", "Setting", "Value")
	fmt.Fprintf(writer, "\n %s\t%s", "–––––––", "–––––––")

	results := env.GetConfig(args)
	for _, config := range results {
		if config.Val == "" {
			config.Val = "–"
		}
		fmt.Fprintf(writer, "\n %s\t%s", strings.ToUpper(config.Key), config.Val)
	}
	fmt.Fprintln(writer, "")
}
