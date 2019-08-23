![Sprawl Logo](assets/logo.png)

![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/eqlabs/sprawl?sort=semver&token=48611096faf7067cc7d8ef9c175f6e7e28f77405)
![CircleCI](https://img.shields.io/circleci/build/github/eqlabs/sprawl/master?token=48611096faf7067cc7d8ef9c175f6e7e28f77405)
[![codecov](https://codecov.io/gh/eqlabs/sprawl/branch/master/graph/badge.svg?token=ms5ajZaWsE)](https://codecov.io/gh/eqlabs/sprawl)
![Matrix](https://img.shields.io/matrix/public:equilibrium.co?server_fqdn=matrix.equilibrium.co)
---

Initial implementation of Sprawl order book protocol in Golang.

# Prerequisites
For developing, preferably a Linux environment with at least Go version 1.11 installed, since the project uses Go Modules. When developing with Windows, the following defaults won't hold:

## Linux
### Create the data directory for Sprawl
```bash
mkdir /var/lib/sprawl
chmod 755 /var/lib/sprawl
```
`sudo` if necessary.

## Windows
### Create an override config file
```bash
cp ./config/default/config.toml ./config.toml
```
The `config.toml` file is ignored in git and it will override every config under `./config`, even under `./config/test`. You need to at least override the database path, since the default directory doesn't exist in Windows.

## More on configuring
The default configuration files reside under `./config`. All the variables there are replaceable by creating a `config.toml` file in project root or defining environment variables with the prefix `SPRAWL_`, for example `SPRAWL_DATABASE_PATH = /var/lib/sprawl/data`

## Generate service code based on the protobuf definition
You only need to do this if something has changed in `./api/service.proto`.
```bash
protoc -I=./api --go_out=plugins=grpc:./api ./api/service.proto
```

## Run all tests (verbose)
```bash
go test -v ./...
```

## Run all tests, see coverage
The following commands generate a code coverage report and open it up in your default web browser.
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```
