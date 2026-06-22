# CeaCtl

Cisco Enterprise API Control

Cisco UCS/MDS API playground written in Go.

## Features

- Cisco MDS API integration
- Cisco UCS XML API integration
- Inventory collection
- Concurrent polling experiments
- Go-Based CLI playground

## Build
```Bash
go mod tidy
go build -o ceactl.exe main.go
```

## Run
```Bash
ceactl.exe mds inventory
ceactl.exe ucsm blades
```

## TODO

- Add MDS Command
- Add UCSM Command
- Add Common Flag (--host, --user, --password, --port, --insecure etc.)
- Credential migrate .env to config.yml with Manual Password type like SSH
- Improve error message for config and device selection
- Add output format option (`--output json`)
- Add result to save a file

### Long-Term (someday)
- Add Local Ollama API function
- Analysis with Local Ollama Model  