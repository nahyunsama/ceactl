# CeaCtl

Cisco Enterprise API Control

Cisco UCS/MDS API playground written in Go.

## Features

- Cisco MDS API integration
- Cisco UCS XML API integration
- Inventory collection
- Concurrent polling experiments
- Go-Based CLI playground

## Project Structure

```
CeaCtl/
├── main.go
├── cmd/
│   ├── root.go
│   ├── mds/
│   │   ├── mds.go            # MDS command entry point, shared flags
│   │   └── show_version.go   # mds version command
│   └── ucsm/
│       ├── ucsm.go           # UCSM command entry point, shared flags
│       └── show_servers.go   # ucsm servers command
└── internal/
    ├── config/
    │   └── config_load.go    # YAML parsing, device selection
    ├── mds/
    │   ├── config/           # MDS-specific configuration
    │   ├── transceiver/      # HTTP client, NX-API requests
    │   ├── receiver/         # JSON response parsing
    │   └── commands/         # Business logic
    └── ucsm/
        ├── config/           # UCSM-specific configuration
        ├── transceiver/      # HTTP client, XML session/requests
        ├── receiver/         # XML response parsing
        └── commands/         # Business logic
```

## Build
```Bash
go mod tidy
go build -o ceactl.exe main.go
```

## Run
```Bash
ceactl.exe mds inventory
ceactl.exe ucsm servers
```

## LLM LOG Analysis Beta
Enable llm_analysis in .config.yaml
```Bash
ceactl.exe mds logs analyze --device <name> --from <date> --to <date>
ceactl.exe mds logs analyze --device <name> --file /path/to/logfile.log
```

```
## TODO

- Add MDS Command
- Add UCSM Command
- Add Common Flag (--host, --user, --password, --port, --insecure etc.)
- Credential migrate .env to config.yml with Manual Password type like SSH
- Add output format option (`--output json`)
- Add result to save a file
- Add MDS Config Command with comfirm (Like Add Zone and zoneset)
- API Level Error Parsing(NX-API code/msg, UCSM errorCode check)
- Add mds logs analyze --around(base time) and --window(before and after) subcommand

### Long-Term (someday)
- Add Local Ollama API function
- Analysis with Local Ollama Model  