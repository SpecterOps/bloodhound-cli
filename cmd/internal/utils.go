package internal

// Various utilities used by other parts of the internal package
// Includes utilities for interacting with the file system

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/adrg/xdg"
)

// HealthIssue is a custom type for storing healthcheck output.
type HealthIssue struct {
	Type    string
	Service string
	Message string
}

type HealthIssues []HealthIssue

func (c HealthIssues) Len() int {
	return len(c)
}

func (c HealthIssues) Less(i, j int) bool {
	return c[i].Service < c[j].Service
}

func (c HealthIssues) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

// GetCwdFromExe gets the current working directory based on "bloodhound-cli" location.
func GetCwdFromExe() string {
	exe, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get path to current executable.")
	}
	return filepath.Dir(exe)
}

// FileExists determines if a given string is a valid filepath.
// Reference: https://golangcode.com/check-if-a-file-exists/
func FileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return !info.IsDir()
}

// DirExists reports whether the specified path exists and is a directory. Returns false if the path does not exist.
func DirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return info.IsDir()
}

// GetDefaultConfigDir returns the default BloodHound config directory path as a `bloodhound` folder inside the current user's data directory.
// Logs a fatal error if the user's config directory cannot be determined.
func GetDefaultConfigDir() string {
	return filepath.Join(xdg.ConfigHome, "bloodhound")
}

// GetBloodHoundDir returns the configured BloodHound config directory path from the environment variable "config_directory".
func GetBloodHoundDir() string {
	return bhEnv.GetString("config_directory")
}

// MakeConfigDir ensures the configured BloodHound config directory exists, creating it if necessary.
// Returns an error if directory creation fails.
func MakeConfigDir() error {
	configDir := GetBloodHoundDir()
	if !DirExists(configDir) {
		log.Printf("The BloodHound config directory you have set, %s, is missing, so attempting to create it.\n", configDir)
		mkErr := os.MkdirAll(configDir, 0777)
		if mkErr != nil {
			return mkErr
		}
		log.Println("Successfully created the BloodHound config directory.")
	}

	return nil
}

// CheckConfigDir checks if the config directory's permissions are at least 0600. This ensures the current user has R/W
// access and BloodHound CLI can function. A more permissive mode won't trigger any errors.
// It returns true if the permissions are sufficient, along with any error encountered during the stat operation.
func CheckConfigDir(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	baselinePerms := os.FileMode(0600)
	mode := info.Mode().Perm()
	return mode&baselinePerms == baselinePerms, nil
}

// GetYamlFilePath joins and returns the directory path of the BloodHound config directory with the Docker Compose YAML file.
// If a user has provided the `-f` or `--file` flag with a string value, the function will return that filepath.
func GetYamlFilePath(override string) string {
	if override != "" {
		log.Printf("Using the override filepath: %s", override)
		fileInfo, err := os.Stat(override)
		if err != nil {
			if os.IsNotExist(err) {
				log.Fatalf("The override path '%s' does not exist.\n", override)
			} else {
				log.Fatalf("There was an error checking the override path '%s': %v\n", override, err)
			}
		}
		if fileInfo.IsDir() {
			log.Fatalf("The provided override path '%s' is a directory instead of a YAML file.\n", override)
		} else {
			if !FileExists(override) {
				log.Fatalf("Override filepaths does not exist: %s", override)
			}
		}
		return override
	}
	return filepath.Join(GetBloodHoundDir(), "docker-compose.yml")
}

// CheckYamlExists verifies that a YAML file exists at the specified path.
// If the file does not exist, it logs a fatal error with instructions for obtaining the required YAML file.
func CheckYamlExists(path string) {
	if !FileExists(path) {
		log.Fatalf(
			"The YAML file %s does not exist! To continue, move your YAML file into the config directory or run "+
				"`./bloodhound-cli check` to download the necessary YAML file.",
			path)
	}
}

// CheckPath returns true if the specified command exists in the system's PATH.
func CheckPath(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// RunBasicCmd executes a given command ("name") with a list of arguments ("args")
// and returns a "string" with the output.
func RunBasicCmd(name string, args []string) (string, error) {
	out, err := exec.Command(name, args...).Output()
	output := string(out[:])
	return output, err
}

// RunCmd executes a given command ("name") with a list of arguments ("args")
// and return stdout and stderr buffers.
func RunCmd(name string, args []string) error {
	// If the command is ``docker`` or ``podman``, prepend ``compose`` to the args
	if name == "docker" || name == "podman" {
		args = append([]string{"compose"}, args...)
	}
	path, err := exec.LookPath(name)
	if err != nil {
		log.Fatalf("`%s` is not installed or not available in the current PATH variable", name)
	}
	exe, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get path to current executable.")
	}
	exePath := filepath.Dir(exe)
	command := exec.Command(path, args...)
	command.Dir = exePath

	stdout, err := command.StdoutPipe()
	if err != nil {
		log.Fatalf("Failed to get stdout pipe for running `%s`", name)
	}
	stderr, err := command.StderrPipe()
	if err != nil {
		log.Fatalf("Failed to get stderr pipe for running `%s`", name)
	}

	stdoutScanner := bufio.NewScanner(stdout)
	stderrScanner := bufio.NewScanner(stderr)
	go func() {
		for stdoutScanner.Scan() {
			fmt.Printf("%s\n", stdoutScanner.Text())
		}
	}()
	go func() {
		for stderrScanner.Scan() {
			fmt.Printf("%s\n", stderrScanner.Text())
		}
	}()
	err = command.Start()
	if err != nil {
		log.Fatalf("Error trying to start `%s`: %v\n", name, err)
	}
	err = command.Wait()
	if err != nil {
		fmt.Printf("[-] Error from `%s`: %v\n", name, err)
		return err
	}
	return nil
}

// Contains checks if a slice of strings ("slice" parameter) contains a given
// string ("search" parameter).
func Contains(slice []string, target string) bool {
	for _, item := range slice {
		if item == target {
			return true
		}
	}
	return false
}

// Silence any output from tests.
// Place `defer quietTests()()` after test declarations.
// Ref: https://stackoverflow.com/a/58720235
func quietTests() func() {
	null, _ := os.Open(os.DevNull)
	sout := os.Stdout
	serr := os.Stderr
	os.Stdout = null
	os.Stderr = null
	log.SetOutput(null)
	return func() {
		defer null.Close()
		os.Stdout = sout
		os.Stderr = serr
		log.SetOutput(os.Stderr)
	}
}

// AskForConfirmation asks the user for confirmation. A user must type in "yes" or "no" and
// then press enter. It has fuzzy matching, so "y", "Y", "yes", "YES", and "Yes" all count as
// confirmations. If the input is not recognized, it will ask again. The function does not return
// until it gets a valid response from the user.
// Original source: https://gist.github.com/r0l1/3dcbb0c8f6cfe9c66ab8008f55f8f28b
func AskForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}

// DownloadFile downloads a file from the specified URL and saves it to the provided filepath.
func DownloadFile(url string, filepath string) error {
	// Create the file to stream the contents
	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Fetch the file to begin the download
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: received status code %d", resp.StatusCode)
	}

	// Write the contents to the file created earlier
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// GetRemoteBloodHoundCliVersion fetches the latest BloodHound CLI version from GitHub's API.
func GetRemoteBloodHoundCliVersion() (string, string, error) {
	var output string

	baseUrl := "https://api.github.com/repos/SpecterOps/bloodhound-cli/releases/latest"
	client := http.Client{Timeout: time.Second * 10}
	resp, err := client.Get(baseUrl)
	if err != nil {
		return "", "", err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("unexpected HTTP status: %d", resp.StatusCode)
	}
	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return "", "", readErr
	}

	var githubJson map[string]interface{}
	jsonErr := json.Unmarshal(body, &githubJson)
	if jsonErr != nil {
		return "", "", jsonErr
	}

	publishedAtRaw, ok := githubJson["published_at"]
	if !ok {
		return "", "", fmt.Errorf("missing 'published_at' in GitHub response")
	}
	publishedAt, ok := publishedAtRaw.(string)
	if !ok {
		return "", "", fmt.Errorf("'published_at' is not a string")
	}
	date, parseErr := time.Parse(time.RFC3339, publishedAt)
	if parseErr != nil {
		output = fmt.Sprintf("BloodHound CLI (published at: %s)", publishedAt)
	} else {
		tagNameRaw, ok := githubJson["tag_name"]
		if !ok {
			return "", "", fmt.Errorf("missing 'tag_name' in GitHub response")
		}
		tagName, ok := tagNameRaw.(string)
		if !ok {
			return "", "", fmt.Errorf("'tag_name' is not a string")
		}
		output = fmt.Sprintf(
			"BloodHound CLI %s (%02d %s %d)",
			tagName, date.Day(), date.Month().String(), date.Year(),
		)
	}

	urlRaw, ok := githubJson["html_url"]
	if !ok {
		return "", "", fmt.Errorf("missing 'html_url' in GitHub response")
	}
	url, ok := urlRaw.(string)
	if !ok {
		return "", "", fmt.Errorf("'html_url' is not a string")
	}
	return output, url, nil
}
