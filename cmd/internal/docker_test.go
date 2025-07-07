package internal

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestEvaluateDockerComposeStatus(t *testing.T) {
	// Mock the BloodHound Docker YAML files
	dataDir := GetDefaultDataDir()
	err := os.MkdirAll(dataDir, 0755)
	assert.NoError(t, err, "Expected directory creation to succeed")

	localMockYaml := filepath.Join(GetDefaultDataDir(), "docker-compose.dev.yml")
	local, localErr := os.Create(localMockYaml)
	prodMockYaml := filepath.Join(GetDefaultDataDir(), "docker-compose.yml")
	prod, prodErr := os.Create(prodMockYaml)
	assert.Equal(t, nil, localErr, "Expected `os.Create()` to return no error")
	assert.Equal(t, nil, prodErr, "Expected `os.Create()` to return no error")
	assert.True(t, FileExists(localMockYaml), "Expected `FileExists()` to return true")
	assert.True(t, FileExists(prodMockYaml), "Expected `FileExists()` to return true")

	defer local.Close()
	defer prod.Close()
}
