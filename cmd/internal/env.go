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
// Defaults are geared towards a development environment.
func setBloodHoundConfigDefaultValues() {
	bhEnv.SetDefault("version", "1")

	// Initial user setup
	bhEnv.SetDefault("default_admin.principal_name", "admin")
	bhEnv.SetDefault("default_admin.bh_default_admin_password", GenerateRandomPassword(32, true))

	// Base config
	bhEnv.SetDefault("bind_addr", "0.0.0.0:8080")
	bhEnv.SetDefault("metrics_port", ":2112")
	bhEnv.SetDefault("root_url", "http://127.0.0.1:8080/")
	bhEnv.SetDefault("work_dir", "/opt/bloodhound/work")
	bhEnv.SetDefault("log_level", "INFO")
	bhEnv.SetDefault("log_path", "bloodhound.log")
	bhEnv.SetDefault("collectors_base_path", "/etc/bloodhound/collectors")

	// TLS config
	bhEnv.SetDefault("tls.cert_file", "")
	bhEnv.SetDefault("tls.key_file", "")

	// Set some helpful aliases for common settings
	bhEnv.RegisterAlias("default_password", "default_admin.bh_default_admin_password")
}

// WriteBloodHoundEnvironmentVariables writes the environment variables to the JSON config file.
func WriteBloodHoundEnvironmentVariables() {
	checkJsonFileExistsAndCreate()
	err := bhEnv.WriteConfig()
	if err != nil {
		log.Fatalf("Error while writing the JSON config file: %s", err)
	}
}

// checkJsonFileExistsAndCreate checks if the JSON file exists and creates it with an empty value, {}, if it doesn't.
func checkJsonFileExistsAndCreate() {
	if !FileExists(filepath.Join(GetCwdFromExe(), "bloodhound.config.json")) {
		file, err := os.Create(filepath.Join(GetCwdFromExe(), "ƒ"))

		if err != nil {
			log.Fatalf("The JSON config file doesn't exist and couldn't be created")
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
	}
}

// ParseBloodHoundEnvironmentVariables attempts to find and open an existing JSON config file or create a new one.
// If a JSON config file is found, load it into the Viper configuration.
// If a JSON config file is not found, create a new one with default values.
// Then write the final file with `WriteBloodHoundEnvironmentVariables()`.
func ParseBloodHoundEnvironmentVariables() {
	setBloodHoundConfigDefaultValues()
	bhEnv.SetConfigName("bloodhound.config.json")
	bhEnv.SetConfigType("json")
	bhEnv.AddConfigPath(GetCwdFromExe())
	bhEnv.AutomaticEnv()
	// Check if expected env file exists
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
			log.Fatalf("Config variable `%s` not found", setting)
		} else {
			values = append(values, Configuration{setting, val})
		}
	}

	sort.Sort(values)

	return values
}

// SetConfig sets the value of the specified key in the JSON config file.
func SetConfig(key string, value string) {
	if strings.ToLower(value) == "true" {
		bhEnv.Set(key, true)
	} else if strings.ToLower(value) == "false" {
		bhEnv.Set(key, false)
	} else {
		bhEnv.Set(key, value)
	}
	WriteBloodHoundEnvironmentVariables()
}
