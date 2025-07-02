package internal

// Various utilities used by other parts of the internal package
// Includes utilities for interacting with the file system

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
		log.Fatalf("Failed to get path to current executable")
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

// DirExists determines if a given string is a valid directory.
// Reference: https://golangcode.com/check-if-a-file-exists/
func DirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return info.IsDir()
}

// GetDefaultHomeDir returns the path for the default BloodHound home directory for initial config creation.
// The initial path will always be a hidden `.BloodHound` directory inside the current user's home directory.
func GetDefaultHomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get user's home directory path to set default home directory: %v", err)
	}
	return filepath.Join(homeDir, ".BloodHound")
}

// GetBloodHoundDir returns the full path configured as the home directory.
func GetBloodHoundDir() string {
	homeDir := bhEnv.GetString("home_directory")
	return filepath.Join(homeDir)
}

// MakeHomeDir checks if the configured home directory exists and creates it if it does not.
func MakeHomeDir() error {
	homeDir := GetBloodHoundDir()
	if !DirExists(homeDir) {
		log.Printf("The configured BloodHound home directory, %s, is missing, so attempting to create it\n", homeDir)
		mkErr := os.MkdirAll(homeDir, 0600)
		if mkErr != nil {
			return mkErr
		}
		log.Println("Successfully created the BloodHound home directory")
	}

	return nil
}

// CheckHomeDir checks if the home directory's permissions are at least 0600. This ensures the current user has full
// access and no other users have access. This is intended as a secure baseline. A more permissive permissions set won't
// trigger any errors.
func CheckHomeDir(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	baselinePerms := os.FileMode(0600)
	mode := info.Mode().Perm()
	return mode&baselinePerms == baselinePerms, nil
}

// DeleteDir deletes the configured home directory and all contents. This is intended as the final step of the
// `uninstall` command.
func DeleteDir(path string) error {
	if DirExists(path) {
		delErr := os.RemoveAll(path)
		if delErr != nil {
			return delErr
		}
	}

	return nil
}

// CheckYamlExists determines if the specified file exists and logs a fatal warning if it does not. It is a wrapper for
// the `FileExists` function and is intended to check YAML files just before executing Docker commands.
func CheckYamlExists(path string) {
	if !FileExists(path) {
		log.Fatalf(
			"The YAML file %s does not exist! To continue, move your YAML file into the home directory or run "+
				"`./bloodhound-cli check` to download the necessary YAML file.",
			path)
	}
}

// CheckPath checks the $PATH environment variable for a given "cmd" and return a "bool"
// indicating if it exists.
func CheckPath(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// RunBasicCmd executes a given command ("name") with a list of arguments ("args")
// and return a "string" with the output.
func RunBasicCmd(name string, args []string) (string, error) {
	out, err := exec.Command(name, args...).Output()
	output := string(out[:])
	return output, err
}

// RunCmd executes a given command ("name") with a list of arguments ("args")
// and return stdout and stderr buffers.
func RunCmd(name string, args []string) error {
	// If the command is ``docker``, prepend ``compose`` to the args
	if name == "docker" {
		args = append([]string{"compose"}, args...)
	}
	path, err := exec.LookPath(name)
	if err != nil {
		log.Fatalf("`%s` is not installed or not available in the current PATH variable", name)
	}
	exe, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get path to current executable")
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
