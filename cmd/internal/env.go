package internal

// Functions for managing the environment variables that control the
// configuration of the BloodHound containers.

import (
	"encoding/json"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Configuration is a custom type for storing configuration values as Key:Val pairs.
type Configuration struct {
	Key string
	Val string
}

// Configurations is a custom type for storing `Configuration` values
type Configurations []Configuration

// Len returns the length of a Configurations struct
func (c Configurations) Len() int {
	return len(c)
}

// Less determines if one Configuration is less than another Configuration
func (c Configurations) Less(i, j int) bool {
	return c[i].Key < c[j].Key
}

// Swap exchanges the position of two Configuration values in a Configurations struct
func (c Configurations) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

// Initialize the environment variables.
var bhEnv = viper.New()

// Set sane defaults for a basic BloodHound deployment.
// setBloodHoundConfigDefaultValues sets default configuration values for BloodHound, including version, admin credentials, server settings, logging, TLS paths, and directory locations. Defaults are intended for development environments.
func setBloodHoundConfigDefaultValues() {
	bhEnv.SetDefault("version", 1)

	// Initial user setup
	bhEnv.SetDefault("default_admin.principal_name", "admin")
	bhEnv.SetDefault("default_admin.password", GenerateRandomPassword(32, true))

	// Base config
	bhEnv.SetDefault("bind_addr", "0.0.0.0:8080")
	bhEnv.SetDefault("metrics_port", ":2112")
	bhEnv.SetDefault("root_url", "http://127.0.0.1:8080")
	bhEnv.SetDefault("work_dir", "/opt/bloodhound/work")
	bhEnv.SetDefault("log_level", "INFO")
	bhEnv.SetDefault("log_path", "bloodhound.log")
	bhEnv.SetDefault("collectors_base_path", "/etc/bloodhound/collectors")
	bhEnv.SetDefault("RecreateDefaultAdmin", "false")

	// TLS config
	bhEnv.SetDefault("tls.cert_file", "")
	bhEnv.SetDefault("tls.key_file", "")

	// Set some helpful aliases for common settings
	bhEnv.RegisterAlias("default_password", "default_admin.password")

	bhEnv.SetDefault("config_directory", GetDefaultConfigDir())
}

// WriteBloodHoundEnvironmentVariables writes the current BloodHound configuration to the JSON config file, ensuring the file exists before writing. Logs a fatal error and exits if writing fails.
func WriteBloodHoundEnvironmentVariables() {
	checkJsonFileExistsAndCreate()
	err := bhEnv.WriteConfig()
	if err != nil {
		log.Fatalf("Error while writing the JSON config file: %s", err)
	}
}

// checkJsonFileExistsAndCreate ensures that the BloodHound JSON configuration file exists in the designated directory
// with proper permissions, creating the file and config directory if necessary. If the file or directory cannot be
// created or permissions are insufficient, the function logs a fatal error and terminates the program.
func checkJsonFileExistsAndCreate() {
	if !FileExists(filepath.Join(GetBloodHoundDir(), "bloodhound.config.json")) {
		configErr := MakeConfigDir()
		if configErr != nil {
			log.Fatalf("Error creating config directory: %s", configErr)
		}

		file, err := os.Create(filepath.Join(GetBloodHoundDir(), "bloodhound.config.json"))

		if err != nil {
			log.Fatalf("The JSON config file doesn't exist and couldn't be created.")
		}

		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Fatalf("Failed to close file: %v", err)
			}
		}(file)

		emptyJSON := make(map[string]interface{})
		encoder := json.NewEncoder(file)
		if err := encoder.Encode(emptyJSON); err != nil {
			log.Fatalf("Failed to write JSON to file: %v", err)
		}
	} else {
		permCheck, permErr := CheckConfigDir(GetBloodHoundDir())
		if permErr != nil {
			log.Fatalf("Error checking the permissions on the config directory: %s", permErr)
		}

		if !permCheck {
			log.Fatalf("The permissions set on the config directory, %s, must be at least allow read and write for the current user (e.g., 0600).", GetBloodHoundDir())
		}
	}
}

// ParseBloodHoundEnvironmentVariables initializes default configuration values, ensures the BloodHound config file and
// directory exist with correct permissions, loads configuration from the JSON file and environment variables, and
// writes the final configuration back to the file. The function terminates the program on critical errors.
func ParseBloodHoundEnvironmentVariables() {
	setBloodHoundConfigDefaultValues()
	bhEnv.SetConfigName("bloodhound.config.json")
	bhEnv.SetConfigType("json")
	bhEnv.AddConfigPath(GetBloodHoundDir())
	bhEnv.AutomaticEnv()
	// Check if the expected JSON file exists
	checkJsonFileExistsAndCreate()
	// Try reading the env file
	if err := bhEnv.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalf("Error while reading in the JSON config file: %s", err)
		} else {
			log.Fatalf("Error while parsing the JSON config file: %s", err)
		}
	}
	WriteBloodHoundEnvironmentVariables()
}

// GetConfigAll retrieves all values from the JSON config configuration file.
func GetConfigAll() []byte {
	configuration := bhEnv.AllSettings()
	configJSON, err := json.MarshalIndent(configuration, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal configuration to JSON: %v", err)
	}

	return configJSON
}

// GetConfig retrieves the specified values from the JSON config file.
func GetConfig(args []string) Configurations {
	var values Configurations
	for i := 0; i < len(args[0:]); i++ {
		setting := strings.ToLower(args[i])
		val := bhEnv.GetString(setting)
		if val == "" {
			log.Fatalf("Config variable `%s` not found.", setting)
		} else {
			values = append(values, Configuration{setting, val})
		}
	}

	sort.Sort(values)

	return values
}

// SetConfig sets the value of the specified key in the JSON config file.
func SetConfig(key string, value string) {
	// We do not support changing the `config_directory` at this time. We can explore that at a later time.
	// Allowing it to be changed will cause the new directory to be created with a blank config file, so we disable the
	// option here to avoid any confusion.
	if strings.ToLower(key) == "config_directory" {
		log.Fatalf("The config directory cannot be changed here, but you can use `--file` a different Docker YAML file to use.")
	}

	if strings.ToLower(value) == "true" {
		bhEnv.Set(key, true)
	} else if strings.ToLower(value) == "false" {
		bhEnv.Set(key, false)
	} else {
		bhEnv.Set(key, value)
	}

	WriteBloodHoundEnvironmentVariables()
}
