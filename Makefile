.PHONY: test, testv, benchmark

protoc:
	protoc --go_out=plugins=grpc:. --cobra_out=plugins=client:. pb/sprawl.proto && protoc -I=./pb --go_out=plugins=grpc:./pb ./pb/sprawl.proto

build: protoc buildwithflags

buildwithflags:
	go build -ldflags "-X main.configPath="

test:
	go test -coverprofile=coverage.out -p 1 ./...

testv:
	go test -coverprofile=coverage.out -p 1 -v ./...

benchmark:
	go test -bench=. -run=^Benchmark ./...

coverage:
	go test -coverprofile=coverage.out -p 1 ./... && go tool cover -html=coverage.out
