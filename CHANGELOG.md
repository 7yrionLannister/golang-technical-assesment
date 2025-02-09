# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

### Changed

### Removed

## [1.3.4] - 2025-02-09

### Added

- v1.3.4 add more functions to database abstraction

## Changed

- v1.3.4 optimize computation of total energy consumption by delegating it to the database

## [1.3.3] - 2025-02-04

### Changed

- v1.3.3 use zap logging framework instead of log/slog for efficiency

## [1.3.2] - 2025-01-26

### Changed

- v1.3.2 fix Docker files

## [1.3.1] - 2025-01-26

### Changed

- v1.3.1 fix multiple meters ids in query params

## [1.3.0] - 2025-01-26

### Added

- v1.3 report daily energy consumption

## [1.2.0] - 2025-01-26

### Added

- v1.2 report weekly energy consumption

## [1.1.0] - 2025-01-26

### Added

- v1.1 report monthly energy consumption
- gin-gonic/gin framework
- swaggo openapi documentation
- gorm model
- import test data from csv file
- .env file configuration
- database model
- database migrations
- README.md template
- CHANGELOG.md template
- Dockerfile template
- docker-compose.yml template
- Initial go.mod