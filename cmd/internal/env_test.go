package internal

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log"
	"path/filepath"
	"testing"
)

// CountConfigProperties returns the number of keys in the JSON configuration file.
func CountConfigProperties() int {
	config := GetConfigAll()
	var configMap map[string]interface{}
	if err := json.Unmarshal(config, &configMap); err != nil {
		log.Fatalf("Failed to unmarshal configuration: %v", err)
	}
	return len(configMap)
}

func TestBloodHoundEnvironmentVariables(t *testing.T) {
	//defer quietTests()()

	// Test parsing values and writing to the JSON config file
	envFile := filepath.Join(GetCwdFromExe(), "bloodhound.config.json")
	ParseBloodHoundEnvironmentVariables()
	assert.True(t, FileExists(envFile), "Expected the JSON file to exist")

	// Test a default value
	assert.Equal(t, bhEnv.Get("default_admin.principal_name"), "admin", "Value of `principal_name` should be `admin`")

	// Test ``GetConfig()``
	format := GetConfig([]string{"collectors_base_path", "default_admin.principal_name"})
	assert.Equal(
		t,
		format,
		Configurations{Configuration{Key: "collectors_base_path", Val: "/etc/bloodhound/collectors"}, Configuration{Key: "default_admin.principal_name", Val: "admin"}},
		"`GetConfig()` should return a Configurations object",
	)
	assert.Equal(t, len(format), 2, "`GetConfig()` with two valid variables should return a two values")

	// Test ``GetConfigAll()``
	assert.Equal(t, 12, CountConfigProperties(), "`GetConfigAll()` should return all values")

	// Test ``SetConfig()``
	SetConfig("log_path", "bhce.log")
	assert.Equal(t, bhEnv.GetString("log_path"), "bhce.log", "New value of `log_path` should be `bhce.log`")
}
