# Golang Technical Assesment - Energy Consumption Report

## Overview
This is a Golang-based microservice designed to report the energy consumption of a set of electricity meters.

## Features
- TODO

## Prerequisites
- Go 1.23.x
- Docker (optional, if using containerization)

## Installation
1. Clone the repository:
    ```sh
    git clone https://github.com/7yrionLannister/golang-technical-assesment.git
    ```
2. Navigate to the project directory:
    ```sh
    cd golang-technical-assesment/app/src
    ```
3. Install dependencies:
    ```sh
    go mod tidy
    ```

## Usage
1. Start the external dependencies (database):
    ```sh
    docker compose up
    ```
2. Run the service:
    ```sh
    go run main.go
    ```
3. Access the service at `http://localhost:8181` (or the configured port).

## Configuration
Configuration options can be set via environment variables:
- `PORT`: The port on which the service runs (default: 8181)
- `DB_HOST`: Database host (and port)
- `DB_USER`: Database user
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name

## Testing
Run the tests using:
```sh
go test -v .\repository\ .\service\
```

## Docker
To build and run the service using Docker:
1. Build the Docker image:
    ```sh
    docker build -t energy-consumption-reporter .
    ```
2. Run the Docker container:
    ```sh
    docker run -p 8181:8181 energy-consumption-reporter
    ```

## Contributing
Contributions are welcome! Please open an issue or submit a pull request.
