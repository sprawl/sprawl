.PHONY: test, benchmark

protoc:
	protoc --go_out=plugins=grpc:. --cobra_out=plugins=client:. pb/sprawl.proto && protoc -I=./pb --go_out=plugins=grpc:./pb ./pb/sprawl.proto

build:
	protoc && go build -ldflags "-X main.configPath="

test:
	go test -p 1 ./...

benchmark:
	go test -bench=. -run=^Benchmark ./...

coverage:
	go test -coverprofile=coverage.out -p 1 ./... && go tool cover -html=coverage.out
