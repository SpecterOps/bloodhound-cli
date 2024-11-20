# BloodHound_CLI

[![Go](https://img.shields.io/badge/Go-1.18-9cf)](.)

![GitHub Release (Latest by Date)](https://img.shields.io/github/v/release/SpecterOps/BloodHound_CLI?label=Latest%20Release)
![GitHub Release Date](https://img.shields.io/github/release-date/SpecterOps/BloodHound_CLI?label=Release%20Date)

Golang code for the `bloodhound-cli` binary in [BloodHound](https://github.com/SpecterOps/BloodHound). This binary provides control for various aspects of BloodHound's configuration.

## Usage

Execute `./bloodhound-cli help` for usage information (see below). More information about BloodHound and how to manage it with `bloodhound-cli` can be found on the [BloodHound Wiki](https://github.com/SpecterOps/BloodHound/wiki/).

## Compilation

Releases are compiled with the following command to set version and build date information:

```bash
go build -ldflags="-s -w -X 'github.com/SpecterOps/BloodHound_CLI/cmd/config.Version=`git describe --tags --abbrev=0`' -X 'github.com/SpecterOps/BloodHound_CLI/cmd/config.BuildDate=`date -u '+%d %b %Y'`'" -o bloodhound-cli main.go
```

The version for rolling releases is set to `rolling`.
