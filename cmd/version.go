package cmd

import (
	"fmt"
	"github.com/SpecterOps/BloodHound_CLI/cmd/config"
	utils "github.com/SpecterOps/BloodHound_CLI/cmd/internal"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Displays BloodHound CLI's version information",
	Long: `Displays BloodHound CLI's version information. The local version information comes from the current binary.
The latest release information is pulled from GitHub's API`,
	RunE: compareCliVersions,
}

// init registers the version command with the root command, enabling the "version" CLI command.
func init() {
	rootCmd.AddCommand(versionCmd)
}

// compareCliVersions collects BloodHound CLI's local and latest stable release version numbers and build dates and then
// prints them to standard output.
func compareCliVersions(cmd *cobra.Command, args []string) error {
	// initialize tabwriter
	writer := new(tabwriter.Writer)
	// Set minwidth, tabwidth, padding, padchar, and flags
	writer.Init(os.Stdout, 8, 8, 1, '\t', 0)

	defer writer.Flush()

	fmt.Println("[+] Fetching latest version information:")

	remoteVersion, htmlUrl, remoteErr := utils.GetRemoteBloodHoundCliVersion()
	if remoteErr != nil {
		return remoteErr
	}
	if len(config.BuildDate) == 0 {
		fmt.Fprintf(writer, "\nLocal Version\tBloodHound CLI %s", config.Version)
	} else {
		fmt.Fprintf(writer, "\nLocal Version\tBloodHound CLI %s (%s)", config.Version, config.BuildDate)
	}
	fmt.Fprintf(writer, "\nLatest Release\t%s", remoteVersion)
	fmt.Fprintf(writer, "Latest Download URL\t%s\n", htmlUrl)

	return nil
}
