.PHONY: protoc

protoc:
	protoc --go_out=plugins=grpc:. --cobra_out=plugins=client:. pb/sprawl.proto && protoc -I=./pb --go_out=plugins=grpc:./pb ./pb/sprawl.proto
