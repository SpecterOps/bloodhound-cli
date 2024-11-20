package config

// Constants and variables used by the BloodHound CLI

var (
	// BloodHound CLI version
	// This gets populated at build time with the following flags:
	//   go build -ldflags="-s -w \
	//   -X 'github.com/SpecterOps/BloodHound_CLI/cmd/config.Version=`git describe --tags --abbrev=0`' \
	//   -X 'github.com/SpecterOps/BloodHound_CLI/cmd/config.BuildDate=`date -u '+%d %b %Y'`'" \
	//   -o bloodhound-cli main.go
	Version     string = "v0.1.0"
	BuildDate   string
	Name        string = "BloodHound CLI"
	DisplayName string = "BloodHound CLI"
	Description string = "A command line interface for BloodHound Community Edition"
)
