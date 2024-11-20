package internal

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestEvaluateDockerComposeStatus(t *testing.T) {
	// Mock the BloodHound Docker YAML files
	localMockYaml := filepath.Join(GetCwdFromExe(), "docker-compose.yml")
	local, localErr := os.Create(localMockYaml)
	assert.Equal(t, nil, localErr, "Expected `os.Create()` to return no error")
	assert.Equal(t, nil, prodErr, "Expected `os.Create()` to return no error")
	assert.True(t, FileExists(localMockYaml), "Expected `FileExists()` to return true")
	assert.True(t, FileExists(prodMockYaml), "Expected `FileExists()` to return true")

	defer local.Close()
	defer prod.Close()

	result := EvaluateDockerComposeStatus()
	assert.NoError(t, result, "Expected `EvaluateDockerComposeStatus()` to return no error")
}
