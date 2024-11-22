package cmd

import (
	"fmt"
	docker "github.com/SpecterOps/BloodHound_CLI/cmd/internal"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"text/tabwriter"
)

// runningCmd represents the running command
var runningCmd = &cobra.Command{
	Use:   "running",
	Short: "Print a list of running BloodHound services",
	Long: `Print a list of running BloodHound services.

If containers are found, the results will include information similar
the information provided by the "docker containers ls" command.`,
	Run: displayRunning,
}

func init() {
	rootCmd.AddCommand(runningCmd)
}

func displayRunning(cmd *cobra.Command, args []string) {
	docker.EvaluateDockerComposeStatus()
	// initialize tabwriter
	writer := new(tabwriter.Writer)
	// Set minwidth, tabwidth, padding, padchar, and flags
	writer.Init(os.Stdout, 8, 8, 1, ' ', 0)

	defer writer.Flush()

	fmt.Println("[+] Collecting list of running BloodHound containers...")

	containers := docker.GetRunning()
	fmt.Printf("[+] Found %d running BloodHound containers\n", len(containers))

	if len(containers) > 0 {
		fmt.Fprintf(writer, "\n %s\t%s\t%s\t%s\t%s", "Name", "Container ID", "Image", "Status", "Ports")
		fmt.Fprintf(writer, "\n %s\t%s\t%s\t%s\t%s", "––––––––––––", "––––––––––––", "––––––––––––", "––––––––––––", "––––––––––––")
		for _, container := range containers {
			var ports []string
			for _, port := range container.Ports {
				var portString string

				if port.IP != "" {
					portString = fmt.Sprintf("%s:", port.IP)
				}
				if port.PrivatePort != 0 {
					portString += fmt.Sprintf("%d", port.PrivatePort)
				}
				if port.PublicPort != 0 {
					portString += fmt.Sprintf(":%d » %d/%s", port.PrivatePort, port.PublicPort, port.Type)
				} else {
					portString += fmt.Sprintf("/%s", port.Type)
				}
				ports = append(ports, portString)
			}
			fmt.Fprintf(writer, "\n %s\t%s\t%s\t%v\t%s", container.Name, container.ID, container.Image, container.Status, strings.Join(ports, ", "))
		}
		fmt.Fprintln(writer, "")
	}
}
