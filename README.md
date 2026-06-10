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
- Add Common Flag (--host, --user, --password, --port, --insecure etc.)
- Add UCSM API Call & Parse function
- Add UCSM Command
- Credential migrate .env to config.yml with Manual Password type like SSH

### Long-Term (someday)
- Add Local Ollama API function
- Analysis with Local Ollama Model  