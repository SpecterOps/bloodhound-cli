package internal

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestBloodHoundEnvironmentVariables(t *testing.T) {
	defer quietTests()()

	// Test parsing values and writing to the .env file
	envFile := filepath.Join(GetCwdFromExe(), ".env")
	ParseBloodHoundEnvironmentVariables()
	assert.True(t, FileExists(envFile), "Expected .env file to exist")

	// Test a default value
	assert.Equal(t, bhEnv.Get("postgres_user"), "bloodhound", "Value of `postgres_user` should be `bloodhound`")

	// Test ``GetConfig()``
	format := GetConfig([]string{"neo4j_allow_upgrade", "neo4j_user"})
	assert.Equal(
		t,
		format,
		Configurations{Configuration{Key: "neo4j_allow_upgrade", Val: "true"}, Configuration{Key: "neo4j_user", Val: "neo4j"}},
		"`GetConfig()` should return a Configurations object",
	)
	assert.Equal(t, len(format), 2, "`GetConfig()` with two valid variables should return a two values")

	// Test ``GetConfigAll()``
	config := GetConfigAll()
	assert.Equal(t, len(config), 14, "`GetConfigAll()` should return all values")

	// Test ``SetConfig()``
	SetConfig("bloodhound_port", "9000")
	assert.Equal(t, bhEnv.GetString("bloodhound_port"), "9000", "New value of `bloodhound_port` should be `9000`")
}
