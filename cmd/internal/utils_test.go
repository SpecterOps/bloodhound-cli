package internal

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestGetCwdFromExe(t *testing.T) {
	cwd := GetCwdFromExe()
	assert.False(t, cwd == "", "Expected `GetCwdFromExe()` to return a non-empty string")
}

func TestCheckPath(t *testing.T) {
	assert.True(t, CheckPath("docker"), "Expected `CheckPath()` to find `docker` in `$PATH`")
}

func TestRunBasicCmd(t *testing.T) {
	defer quietTests()()
	_, err := RunBasicCmd(dockerCmd, []string{"--version"})
	assert.Equal(t, nil, err, "Expected `RunBasicCmd()` to return no error")
}

func TestRunCmd(t *testing.T) {
	defer quietTests()()
	err := RunCmd(dockerCmd, []string{"--version"})
	assert.Equal(t, nil, err, "Expected `RunCmd()` to return no error")
}

func TestContains(t *testing.T) {
	assert.True(t, Contains([]string{"a", "b", "c"}, "b"), "Expected `Contains()` to return true")
	assert.False(t, Contains([]string{"a", "b", "c"}, "d"), "Expected `Contains()` to return false")
}

func TestGetRemoteGhostwriterCliVersion(t *testing.T) {
	// Test reading the version data from GitHub's API
	version, _, err := GetRemoteBloodHoundCliVersion()
	assert.Equal(t, nil, err, "Expected `GetRemoteBloodHoundCliVersion()` to return no error")
	assert.True(
		t,
		strings.Contains(version, "BloodHound CLI v"),
		"Expected `GetRemoteBloodHoundCliVersion()` to return a string containing `BloodHound CLI CLI v...`",
	)
}
