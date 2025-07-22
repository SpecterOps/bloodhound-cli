# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.8] - 2025-07-22

### Added

* Added alias commands for `containers up` and `containers down` to make it easier to run `up` and `down`

### Changed

* The `version` command now pulls the latest stable release version information for comparison and provides a download link

## [0.1.7] - 2025-7-9

### Added

* Added support for a dedicated config directory to act as the configuration home for the JSON configuration file and default Docker YAML files
  * The directory is the user's XDG config home directory and `bloodhound`
    * i.e., the equivalent of `~/.config/bloodhound` on Unix, \
      `~/Library/Application Support/bloodhound` on macOS, and \
      `%LOCALAPPDATA%\bloodhound` on Windows
    * We use a lowercase `bloodhound` to match the directory used by older installations of BloodHound
  * You can place BloodHound CLI anywhere and run it from any location, and it will always look in the config directory for the JSON and default YAML files
  * The CLI creates the directory with a `0777` permissions mask so it is accessible to all BloodHound users in multi-user environments
  * The permissions follow your [umask](https://man7.org/linux/man-pages/man2/umask.2.html), so the typical user mask of `0022` will set the permissions to `0755`
* Added a `config_directory` value to the JSON configuration file to control the config directory path
  * Changing this path will change where BloodHound CLI looks for the Docker YAML files
  * BloodHound CLI will continue to look in the default location for the JSON config file 
* Added checks that ensure the configured directory will work as expected every time BloodHound CLI runs
  * The first check ensures the directory exists and creates the directory if it does not
  * The second check ensures the config directory has proper permissions that will allow BloodHound CLI to read and write
* Added a `-f` or `--file` flag to override the location of the YAML file to use for Docker
  * Providing a file path will override where BloodHound CLI looks for the YAML file
  * e.g., `./bloodhound-cli -f /Users/Mable/BloodHound/custom-docker-compose.yml containers up`

### Changed

* Every command that runs a Docker command will now ensure the required YAML file exists before proceeding

## [0.1.6] - 2025-4-23

### Added

* Added a `check` command to check for necessary Docker and Docker Compose commands and the YAML files

### Changed

* Updated golang.org/x/net

### Fixed

* Fixed YAML files being downloaded to your current working directory instead of the CLI binary's directory

## [0.1.5] - 2025-3-25

### Changed

* Changed releases to drop the release tag form the asset filenames to make it easier to grab the latest binaries
* Updated golang.org/x/net

### [0.1.4] - 2025-1-31

### Added

* Added an `update` command to pull the latest BloodHound images
* Added a `resetpwd` command to recreate the default admin account if access is lost
  * This requires BloodHound v7.1.0

## [0.1.3] - 2025-1-31

### Added

* Added a `--volumes` flag to the `containers down` command that deletes the data volumes when the containers come down
* Added an `uninstall` command that removes the BloodHound environment by deleting containers, images, and volume data

## [0.1.2] - 2025-1-22

### Fixed

* Fixed `install` output not showing the initial password in the output

## [0.1.1] - 2025-1-21

### Fixed

* Fixed setting the default password for the `install` command

### Added

* Initial commit & release

## [0.1.0] - 2024-11-20

### Added

* Initial commit & release

### Changed

* N/A

### Deprecated

* N/A

### Removed

* N/A

### Fixed

* N/A

### Security

* N/A
