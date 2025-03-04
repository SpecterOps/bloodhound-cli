# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
