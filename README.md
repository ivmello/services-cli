# Services CLI - A CLI for managing multiple services

Application for starting multiple services at once.

## Requirements
Go 1.17 or higher.

## Build

- Create file `services.yaml` in the root directory of the project.
- Define the services you want to start and respective commands and paths.
- Run build command:

```bash
GOOS=linux go build -o services-cli ./cmd/main.go
```

In this command you can change the `GOOS` to your OS and the name of the binary (defined with flag `-o`).

## Usage

```bash
./services-cli -s=service1,service2,service3
```