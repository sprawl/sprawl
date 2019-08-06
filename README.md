![Sprawl Logo](assets/logo.png)

Initial implementation of Sprawl order book protocol in Golang. Developing...

## Generate service code based on the protobuf definition
```protoc -I=./api --go_out=plugins=grpc:./api ./api/service.proto```
