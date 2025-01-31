package internal

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"log"
	"os"
	"path/filepath"
)

// Vars for tracking the list of BloodHound images
// Used for filtering the list of containers returned by the Docker client
var (
	prodImages = []string{
		"bhce_bloodhound", "bhce_neo4j", "bhce_postgres",
	}
	devImages = []string{
		"bhce_bloodhound", "bhce_neo4j", "bhce_postgres",
	}
	// Default root command for Docker commands
	dockerCmd = "docker"
	// URLs for the BloodHound compose files
	devYaml  = "docker-compose.dev.yml"
	prodYaml = "docker-compose.yml"
	devUrl   = "https://raw.githubusercontent.com/SpecterOps/BloodHound_CLI/refs/heads/main/docker-compose.dev.yml"
	prodUrl  = "https://raw.githubusercontent.com/SpecterOps/BloodHound_CLI/refs/heads/main/docker-compose.yml"
	loginUri = "/ui/login"
)

// Container is a custom type for storing container information similar to output from "docker containers ls".
type Container struct {
	ID     string
	Image  string
	Status string
	Ports  []types.Port
	Name   string
}

// Containers is a collection of Container structs
type Containers []Container

// Len returns the length of a Containers struct
func (c Containers) Len() int {
	return len(c)
}

// Less determines if one Container is less than another Container
func (c Containers) Less(i, j int) bool {
	return c[i].Image < c[j].Image
}

// Swap exchanges the position of two Container values in a Containers struct
func (c Containers) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

// EvaluateDockerComposeStatus determines if the host has the "docker compose" plugin or the "docker compose"
// script installed and set the global `dockerCmd` variable.
func EvaluateDockerComposeStatus(install ...bool) error {
	isInstall := false
	if len(install) > 0 {
		isInstall = install[0]
	}

	fmt.Println("[+] Checking the status of Docker and the Compose plugin...")
	// Check for ``docker`` first because it's required for everything to come
	dockerExists := CheckPath("docker")
	if !dockerExists {
		log.Fatalln("Docker is not installed on this system, so please install Docker and try again")
	}

	// Check if the Docker Engine is running
	_, engineErr := RunBasicCmd("docker", []string{"info"})
	if engineErr != nil {
		log.Fatalln("Docker is installed on this system, but the daemon is not running")
	}

	// Check for the ``compose`` plugin as our first choice
	_, composeErr := RunBasicCmd("docker", []string{"compose", "version"})
	if composeErr != nil {
		fmt.Println("[+] The `compose` plugin is not installed, so we'll try the deprecated `docker-compose` script")
		composeScriptExists := CheckPath("docker-compose")
		if composeScriptExists {
			fmt.Println("[+] The `docker-compose` script is installed, so we'll use that instead")
			dockerCmd = "docker-compose"
		} else {
			fmt.Println("[+] The `docker-compose` script is also not installed or not in the PATH")
			log.Fatalln("Docker Compose is not installed, so please install it and try again: https://docs.docker.com/compose/install/")
		}
	}

	// Bail out if we're not in the same directory as the YAML files
	// Otherwise, we'll get a confusing error message from the `compose` plugin
	if !isInstall {
		if !FileExists(filepath.Join(GetCwdFromExe(), prodYaml)) || !FileExists(filepath.Join(GetCwdFromExe(), devYaml)) {
			log.Fatalln("BloodHound CLI must be run in the same directory as the `docker-compose.yml` and `docker-compose.dev.yml` files")
		}
	}

	return nil
}

// RunDockerComposeInstall executes the "docker compose" commands for a first-time installation with
// the specified YAML file ("yaml" parameter).
func RunDockerComposeInstall(yaml string) {
	// If the YAML files don't exist, download them from the BloodHound repo
	if FileExists(filepath.Join(GetCwdFromExe(), prodYaml)) {
		c := AskForConfirmation("[*] A production YAML file already exists in the current directory. Do you want to overwrite it and continue with the install?")
		if !c {
			os.Exit(0)
		}
	}
	fmt.Printf("[+] Downloading the production YAML file from %s...\n", prodUrl)
	prodDownloadErr := DownloadFile(prodUrl, prodYaml)
	if prodDownloadErr != nil {
		log.Fatalf("Error trying to download the production YAML file: %v\n", prodDownloadErr)
	}

	if FileExists(filepath.Join(GetCwdFromExe(), devYaml)) {
		c := AskForConfirmation("[*] A development YAML file already exists in the current directory. Do you want to overwrite it and continue with the install?")
		if !c {
			os.Exit(0)
		}
	}
	fmt.Printf("[+] Downloading the development YAML file from %s...\n", devUrl)
	devDownloadErr := DownloadFile(devUrl, devYaml)
	if devDownloadErr != nil {
		log.Fatalf("Error trying to download the development YAML file: %v\n", devDownloadErr)
	}

	buildErr := RunCmd(dockerCmd, []string{"-f", yaml, "pull"})
	if buildErr != nil {
		log.Fatalf("Error trying to build with %s: %v\n", yaml, buildErr)
	}
	upErr := RunCmd(dockerCmd, []string{"-f", yaml, "up", "-d"})
	if upErr != nil {
		log.Fatalf("Error trying to bring up environment with %s: %v\n", yaml, upErr)
	}
	fmt.Println("[+] BloodHound is ready to go!")
	fmt.Printf("[+] You can log in as `%s` with this password: %s\n", bhEnv.GetString("default_admin.principal_name"), bhEnv.GetString("default_admin.password"))
	fmt.Println("[+] You can get your admin password by running: bloodhound-cli config get default_password")
	fmt.Printf("[+] You can access the BloodHound UI at: %s%s\n", bhEnv.GetString("root_url"), loginUri)
}

// RunDockerComposeUninstall executes the "docker compose" commands to bring down containers and remove containers,
// images, and volumes with the specified YAML file ("yaml" parameter).
func RunDockerComposeUninstall(yaml string) {
	c := AskForConfirmation("[!] This command removes all containers, images, and volume data. Are you sure you want to uninstall?")
	if !c {
		os.Exit(0)
	}
	uninstallErr := RunCmd(dockerCmd, []string{"-f", yaml, "down", "--rmi", "all", "-v", "--remove-orphans"})
	if uninstallErr != nil {
		log.Fatalf("Error trying to uninstall with %s: %v\n", yaml, uninstallErr)
	}
	fmt.Println("[+] Uninstall was successful. You can re-install with `./bloodhound-cli install`.")
}

// RunDockerComposeUpgrade executes the "docker compose" commands for re-building or upgrading an
// installation with the specified YAML file ("yaml" parameter).
func RunDockerComposeUpgrade(yaml string) {
	fmt.Printf("[+] Running `%s` commands to build containers with %s...\n", dockerCmd, yaml)
	downErr := RunCmd(dockerCmd, []string{"-f", yaml, "down"})
	if downErr != nil {
		log.Fatalf("Error trying to bring down any running containers with %s: %v\n", yaml, downErr)
	}
	buildErr := RunCmd(dockerCmd, []string{"-f", yaml, "build"})
	if buildErr != nil {
		log.Fatalf("Error trying to build with %s: %v\n", yaml, buildErr)
	}
	upErr := RunCmd(dockerCmd, []string{"-f", yaml, "up", "-d"})
	if upErr != nil {
		log.Fatalf("Error trying to bring up environment with %s: %v\n", yaml, upErr)
	}
	fmt.Println("[+] All containers have been built!")
}

// RunDockerComposeStart executes the "docker compose" commands to start the environment with
// the specified YAML file ("yaml" parameter).
func RunDockerComposeStart(yaml string) {
	fmt.Printf("[+] Running `%s` to restart containers with %s...\n", dockerCmd, yaml)
	startErr := RunCmd(dockerCmd, []string{"-f", yaml, "start"})
	if startErr != nil {
		log.Fatalf("Error trying to restart the containers with %s: %v\n", yaml, startErr)
	}
}

// RunDockerComposeStop executes the "docker compose" commands to stop all services in the environment with
// the specified YAML file ("yaml" parameter).
func RunDockerComposeStop(yaml string) {
	fmt.Printf("[+] Running `%s` to stop services with %s...\n", dockerCmd, yaml)
	stopErr := RunCmd(dockerCmd, []string{"-f", yaml, "stop"})
	if stopErr != nil {
		log.Fatalf("Error trying to stop services with %s: %v\n", yaml, stopErr)
	}
}

// RunDockerComposeRestart executes the "docker compose" commands to restart the environment with
// the specified YAML file ("yaml" parameter).
func RunDockerComposeRestart(yaml string) {
	fmt.Printf("[+] Running `%s` to restart containers with %s...\n", dockerCmd, yaml)
	startErr := RunCmd(dockerCmd, []string{"-f", yaml, "restart"})
	if startErr != nil {
		log.Fatalf("Error trying to restart the containers with %s: %v\n", yaml, startErr)
	}
}

// RunDockerComposeUp executes the "docker compose" commands to bring up the environment with
// the specified YAML file ("yaml" parameter).
func RunDockerComposeUp(yaml string) {
	fmt.Printf("[+] Running `%s` to bring up the containers with %s...\n", dockerCmd, yaml)
	upErr := RunCmd(dockerCmd, []string{"-f", yaml, "up", "-d"})
	if upErr != nil {
		log.Fatalf("Error trying to bring up the containers with %s: %v\n", yaml, upErr)
	}
}

// RunDockerComposeDown executes the "docker compose" commands to bring down the environment with
// the specified YAML file ("yaml" parameter).
func RunDockerComposeDown(yaml string, volumes bool) {
	fmt.Printf("[+] Running `%s` to bring down the containers with %s...\n", dockerCmd, yaml)
	args := []string{"-f", yaml, "down"}
	if volumes {
		args = append(args, "--volumes")
	}
	downErr := RunCmd(dockerCmd, args)
	if downErr != nil {
		log.Fatalf("Error trying to bring down the containers with %s: %v\n", yaml, downErr)
	}
}

// FetchLogs fetches logs from the container with the specified "name" label ("containerName" parameter).
func FetchLogs(containerName string, lines string) []string {
	var logs []string
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Failed to get client in logs: %v", err)
	}
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		log.Fatalf("Failed to get container list: %v", err)
	}
	if len(containers) > 0 {
		for _, container := range containers {
			if container.Labels["name"] == containerName || containerName == "all" || container.Labels["name"] == "bhce_"+containerName {
				logs = append(logs, fmt.Sprintf("\n*** Logs for `%s` ***\n\n", container.Labels["name"]))
				reader, err := cli.ContainerLogs(context.Background(), container.ID, types.ContainerLogsOptions{
					ShowStdout: true,
					ShowStderr: true,
					Tail:       lines,
				})
				if err != nil {
					log.Fatalf("Failed to get container logs: %v", err)
				}
				defer reader.Close()
				// Reference: https://medium.com/@dhanushgopinath/reading-docker-container-logs-with-golang-docker-engine-api-702233fac044
				p := make([]byte, 8)
				_, err = reader.Read(p)
				for err == nil {
					content := make([]byte, binary.BigEndian.Uint32(p[4:]))
					reader.Read(content)
					logs = append(logs, string(content))
					_, err = reader.Read(p)
				}
			}
		}

		if len(logs) == 0 {
			logs = append(logs, fmt.Sprintf("\n*** No logs found for requested container '%s' ***\n", containerName))
		}
	} else {
		fmt.Println("Failed to find that container")
	}
	return logs
}

// GetRunning determines if the container with the specified "name" label ("containerName" parameter) is running.
func GetRunning() Containers {
	var running Containers

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Failed to get client connection to Docker: %v", err)
	}
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
		All: false,
	})
	if err != nil {
		log.Fatalf("Failed to get container list from Docker: %v", err)
	}
	if len(containers) > 0 {
		for _, container := range containers {
			if Contains(devImages, container.Labels["name"]) || Contains(prodImages, container.Labels["name"]) {
				running = append(running, Container{
					container.ID, container.Image, container.Status, container.Ports, container.Labels["name"],
				})
			}
		}
	}

	return running
}
