# BloodHound CLI

![Go](https://img.shields.io/github/go-mod/go-version/SpecterOps/bloodhound-cli?color=50B071)

![GitHub Release (Latest by Date)](https://img.shields.io/github/v/release/SpecterOps/bloodhound-cli?label=Latest%20Release&color=E61616
)
![GitHub Release Date](https://img.shields.io/github/release-date/SpecterOps/bloodhound-cli?label=Release%20Date&color=E1E2EF)

![BHCLI.png](BHCLI.png)

Golang code for the `bloodhound-cli` binary in [BloodHound](https://github.com/SpecterOps/BloodHound). This binary provides control for various aspects of BloodHound's configuration.

## Usage

Execute `./bloodhound-cli help` for usage information (see below). 

More information about BloodHound and how to manage it with `bloodhound-cli` can be found on the [BloodHound Community Edition Quickstart Guide](https://bloodhound.specterops.io/get-started/quickstart/community-edition-quickstart), which is part of the [BloodHound documentation](https://bloodhound.specterops.io/home).

## Compilation

Releases are compiled with the following command to set version and build date information:

```bash
go build -ldflags="-s -w -X 'github.com/SpecterOps/BloodHound_CLI/cmd/config.Version=`git describe --tags --abbrev=0`' -X 'github.com/SpecterOps/BloodHound_CLI/cmd/config.BuildDate=`date -u '+%d %b %Y'`'" -o bloodhound-cli main.go
```

The version for rolling releases is set to `rolling`.
