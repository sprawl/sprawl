![Sprawl Logo](assets/logo.png)

Initial implementation of Sprawl order book protocol in Golang.

# Prerequisites
For developing, a Linux environment with at least Go version 1.11 installed, since the project uses Go Modules.

# Development

## Generate service code based on the protobuf definition
```protoc -I=./api --go_out=plugins=grpc:./api ./api/service.proto```

## Run all tests (verbose)
```go test -v ./...```

## Run all tests, see coverage
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```
