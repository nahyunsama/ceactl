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

