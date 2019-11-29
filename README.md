![Sprawl Logo](assets/logo.png)

[![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/sprawl/sprawl?sort=semver&token=48611096faf7067cc7d8ef9c175f6e7e28f77405)](https://github.com/sprawl/sprawl)
[![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/eqlabs/sprawl)](https://cloud.docker.com/u/eqlabs/repository/docker/eqlabs/sprawl)
[![CircleCI](https://img.shields.io/circleci/build/github/sprawl/sprawl/master?token=48611096faf7067cc7d8ef9c175f6e7e28f77405)](https://circleci.com/gh/sprawl/sprawl/tree/master)
[![codecov](https://codecov.io/gh/sprawl/sprawl/branch/master/graph/badge.svg)](https://codecov.io/gh/sprawl/sprawl)
[![GoDoc](https://godoc.org/github.com/eqlabs/sprawl?status.svg)](https://godoc.org/github.com/sprawl/sprawl)
[![Gitter](https://img.shields.io/gitter/room/sprawl/sprawl)](https://gitter.im/sprawl/sprawl)
---

Sprawl is a peer-to-peer protocol and network for creating distributed marketplaces. It uses [KademliaDHT](https://github.com/libp2p/go-libp2p-kad-dht) for peer discovery, [Libp2p](https://github.com/libp2p) for networking and [LevelDB](https://github.com/google/leveldb) for storage. Sprawl can be combined with any settlement system and security model.

**There is a test node running the debug pinger in the following IP:** `157.245.171.225` This means that if you are running the debug pinger yourself and your node successfully connects to it, you should start receiving _Hello world!_ orders. :)

[Support on Gitter!](https://gitter.im/eqlabs/sprawl)

# Building with Sprawl
## API operations
Sprawl uses protocol buffers for its messaging. The API is described in `./pb/sprawl.proto`

```protobuf
service OrderHandler {
	rpc Create (CreateRequest) returns (CreateResponse);
	rpc Delete (OrderSpecificRequest) returns (GenericResponse);
	rpc Lock (OrderSpecificRequest) returns (GenericResponse);
	rpc Unlock (OrderSpecificRequest) returns (GenericResponse);
	rpc GetOrder (OrderSpecificRequest) returns (Order);
	rpc GetAllOrders (Empty) returns (OrderListResponse);
}

service ChannelHandler {
	rpc Join (JoinRequest) returns (JoinResponse);
	rpc Leave (ChannelSpecificRequest) returns (GenericResponse);
	rpc GetChannel (ChannelSpecificRequest) returns (Channel);
	rpc GetAllChannels (Empty) returns (ChannelListResponse);
}
```

## Configuration options
By default, Sprawl runs on default config which is located under `./config/default/`. You can override these configuration options _during development_ by either creating a config file "config.toml" under root, like `./config.toml`, or _in production_ by using environment variables:

| **Variable**                          | **Description**                                                                                        | **Default**            |
| ------------------------------------- | ------------------------------------------------------------------------------------------------------ | ---------------------- |
| `SPRAWL_RPC_PORT`                     | The gRPC API port                                                                                      | 1337                   |
| `SPRAWL_DATABASE_PATH`                | The folder that LevelDB will use to save its data                                                      | "/var/lib/sprawl/data" |
| `SPRAWL_P2P_DEBUG`                    | Pinger that pushes an order into "testChannel" every minute                                            | false                  |
| `SPRAWL_P2P_ENABLENATPORTMAP` | Enable NAT port mapping on nodes that are behind a firewall. Not compatible with Docker.               | true                  |
| `SPRAWL_P2P_EXTERNALIP` | A public IP to publish for other Sprawl nodes to connect to               | ""                  |
| `SPRAWL_P2P_PORT` | libp2p listen port. Constructs a multiaddress together with EXTERNALIP               | 4001                  |
| `SPRAWL_ERRORS_ENABLESTACKTRACE` | Enable stack trace on error messages               | false                  |

## Running a node
This is the easiest way to run Sprawl. If you only need the default functionality of sending and receiving orders, without any additional fields or any of that sort, this is the recommended way, since you don't need to be informed of Sprawl's internals. It should just work. If it doesn't, create an issue or hit us up on Gitter! :D

```bash
go run main.go
```
OR
```bash
# Build a development version which assumes that it's ran inside the repo with all config files
go build && ./sprawl
# You can also build a binary that doesn't assume any configuration files,
# but in this case you MUST use environment variables
go build -ldflags "-X main.configPath=" && ./sprawl
```
OR
```bash
docker run -it eqlabs/sprawl -p 1337:1337
```

This spawns a Sprawl node with the default configuration. (More information on configuration options at "More on configuring" including environment variables.)

The node then connects to the IPFS bootstrap peers, fetching their DHT routing tables, announcing itself as a part of the Sprawl network.

Different Sprawl nodes should connect to each other using the DHT on the network and open pubsub connections between the channels they're subscribed to. They will then synchronize between each other exchanging `CREATE`, `DELETE`, `LOCK` and `UNLOCK` operations on orders, persisting the state locally on LevelDB.

You can use your or any Sprawl node that's accessible to you with `sprawl-cli`. Documentation on the cli tool is kept separate from this repository. We'd be happy to see you develop your own tools using the gRPC/JSON API of Sprawl!

## Using Sprawl as a library
You can also build your own applications on top of Sprawl using the packages directly in Go. Best way to get a grasp on how this could be done is to check out `./app/app.go` since it's the default application definition which runs a Sprawl node.

Under `./interfaces` you can find the interface definitions that need to be fulfilled. If you want to use just a few packages from or customize Sprawl, you can do it. For example, if you want to replace LevelDB with a different database, you need to program the methods defined in `./interfaces/Storage.go` to fit your specific database, and plug it in the app.

We aim to continuously expand the ways you can make plugins on top of Sprawl.

# Developing Sprawl
## Prerequisites
For developing, preferably a Linux environment with at least Go version 1.11 installed, since the project uses Go Modules. When developing with Windows, the following defaults won't hold:

### Linux
#### Create the data directory for Sprawl
```bash
mkdir /var/lib/sprawl
chmod 755 /var/lib/sprawl
```
`sudo` if necessary.

### Windows
#### Create an override config file
```bash
cp ./config/default/config.toml ./config.toml
```
The `config.toml` file is ignored in git and it will override every config under `./config`, even under `./config/test`. You need to at least override the database path, since the default directory doesn't exist in Windows.

### More on configuring
The default configuration files reside under `./config`. All the variables there are replaceable by creating a `config.toml` file in project root or defining environment variables with the prefix `SPRAWL_`, for example `SPRAWL_DATABASE_PATH = /var/lib/sprawl/data`

### Generate service code based on the protobuf definition
You only need to do this if something has changed in `./pb/sprawl.proto`.
```bash
make protoc

OR

protoc --go_out=plugins=grpc:. --cobra_out=plugins=client:. pb/sprawl.proto && protoc -I=./pb --go_out=plugins=grpc:./pb ./pb/sprawl.proto
```

### Run all tests
```bash
go test -p 1 ./...
```

### Run all tests, see coverage
The following commands generate a code coverage report and open it up in your default web browser.
```bash
go test -coverprofile=coverage.out -p 1 ./...
go tool cover -html=coverage.out
```
