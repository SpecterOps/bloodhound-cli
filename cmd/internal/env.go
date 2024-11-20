package internal

// Functions for managing the environment variables that control the
// configuration of the BloodHound containers.

import (
	"fmt"
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
	bhEnv.SetDefault("bloodhound_tag", "latest")

	// Postgres auth configuration
	bhEnv.SetDefault("postgres_user", "bloodhound")
	bhEnv.SetDefault("postgres_password", GenerateRandomPassword(32, false))
	bhEnv.SetDefault("postgres_db", "bloodhound")
	bhEnv.SetDefault("postgres_db_host", "app-db")

	// Auth string for Neo4j credentials
	bhEnv.SetDefault("neo4j_user", "neo4j")
	bhEnv.SetDefault("neo4j_secret", GenerateRandomPassword(32, false))
	bhEnv.SetDefault("neo4j_host", "graph-db:7687/")

	// Allow upgrades of neo4j data (useful when importing external data)
	bhEnv.SetDefault("neo4j_allow_upgrade", true)

	// Django settings// Port forward information
	bhEnv.SetDefault("bloodhound_host", "127.0.0.1")
	bhEnv.SetDefault("bloodhound_port", 8080)
	bhEnv.SetDefault("postgres_port", 5432)
	bhEnv.SetDefault("neo4j_db_port", 7687)
	bhEnv.SetDefault("neo4j_web_port", "7474")
}

// WriteBloodHoundEnvironmentVariables writes the environment variables to the .env file.
func WriteBloodHoundEnvironmentVariables() {
	c := bhEnv.AllSettings()
	// To make it easier to read and look at, get all the keys, sort them, and display variables in order
	keys := make([]string, 0, len(c))
	for k := range c {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	f, err := os.Create(filepath.Join(GetCwdFromExe(), ".env"))
	if err != nil {
		log.Fatalf("Error writing out environment!\n%v", err)
	}
	defer f.Close()
	for _, key := range keys {
		if len(bhEnv.GetString(key)) == 0 {
			_, err = f.WriteString(fmt.Sprintf("%s=\n", strings.ToUpper(key)))
		} else {
			_, err = f.WriteString(fmt.Sprintf("%s='%s'\n", strings.ToUpper(key), bhEnv.GetString(key)))
		}

		if err != nil {
			log.Fatalf("Failed to write out environment!\n%v", err)
		}
	}
}

// ParseBloodHoundEnvironmentVariables attempts to find and open an existing .env file or create a new one.
// If an .env file is found, load it into the Viper configuration.
// If an .env file is not found, create a new one with default values.
// Then write the final file with `WriteBloodHoundEnvironmentVariables()`.
func ParseBloodHoundEnvironmentVariables() {
	setBloodHoundConfigDefaultValues()
	bhEnv.SetConfigName(".env")
	bhEnv.SetConfigType("env")
	bhEnv.AddConfigPath(GetCwdFromExe())
	bhEnv.AutomaticEnv()
	// Check if expected env file exists
	if !FileExists(filepath.Join(GetCwdFromExe(), ".env")) {
		_, err := os.Create(filepath.Join(GetCwdFromExe(), ".env"))
		if err != nil {
			log.Fatalf("The .env doesn't exist and couldn't be created")
		}
	}
	// Try reading the env file
	if err := bhEnv.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalf("Error while reading in .env file: %s", err)
		} else {
			log.Fatalf("Error while parsing .env file: %s", err)
		}
	}
	WriteBloodHoundEnvironmentVariables()
}

// SetProductionMode updates the environment variables to switch to production mode.
func SetProductionMode() {
	bhEnv.Set("hasura_graphql_dev_mode", false)
	bhEnv.Set("django_secure_ssl_redirect", true)
	bhEnv.Set("django_settings_module", "config.settings.production")
	bhEnv.Set("django_csrf_cookie_secure", true)
	bhEnv.Set("django_session_cookie_secure", true)
	WriteBloodHoundEnvironmentVariables()
}

// GetConfigAll retrieves all values from the .env configuration file.
func GetConfigAll() Configurations {
	c := bhEnv.AllSettings()
	keys := make([]string, 0, len(c))
	for k := range c {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var values Configurations
	for _, key := range keys {
		val := bhEnv.GetString(key)
		values = append(values, Configuration{strings.ToUpper(key), val})
	}

	sort.Sort(values)

	return values
}

// GetConfig retrieves the specified values from the .env file.
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

// SetConfig sets the value of the specified key in the .env file.
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
