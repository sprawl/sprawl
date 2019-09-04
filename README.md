![Sprawl Logo](assets/logo.png)

![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/eqlabs/sprawl?sort=semver&token=48611096faf7067cc7d8ef9c175f6e7e28f77405)
![CircleCI](https://img.shields.io/circleci/build/github/eqlabs/sprawl/master?token=48611096faf7067cc7d8ef9c175f6e7e28f77405)
[![codecov](https://codecov.io/gh/eqlabs/sprawl/branch/master/graph/badge.svg?token=ms5ajZaWsE)](https://codecov.io/gh/eqlabs/sprawl)
![Matrix](https://img.shields.io/matrix/public:equilibrium.co?server_fqdn=matrix.equilibrium.co)
---

Sprawl is a distributed order book protocol and network. Its purpose is to bring together buyers and sellers of any kind of asset, where a viable trade execution mechanism exists, in a single global liquidity pool.

Support on Matrix in `#public:equilibrium.co`

# Running a node
```bash
go run main.go
```
OR
```bash
go build && ./sprawl
```
OR
```bash
docker run -it eqlabs/sprawl -p 1337:1337
```
This spawns a Sprawl node with the default configuration. (More information on configuration options at "More on configuring" including environment variables.)

The node then connects to the IPFS bootstrap peers, fetching their DHT routing tables, announcing itself as a part of the Sprawl network.

Different Sprawl nodes should connect to each other using the DHT on the network and open pubsub connections between the channels they're subscribed to. They will then synchronize between each other exchanging `CREATE`, `DELETE`, `LOCK` and `UNLOCK` operations on orders, persisting the state locally on LevelDB.

You can use your or any Sprawl node that's accessible to you with `sprawl-cli`. Documentation on the cli tool is kept separate from this repository. We'd be happy to see you develop your own tools using the gRPC/JSON API of Sprawl!

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
You only need to do this if something has changed in `./pb/sprawl.proto`.
```bash
make protoc

OR

protoc --go_out=plugins=grpc:. --cobra_out=plugins=client:. pb/sprawl.proto && protoc -I=./pb --go_out=plugins=grpc:./pb ./pb/sprawl.proto
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
