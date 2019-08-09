FROM golang:latest as build

RUN mkdir /app 
COPY . /app/ 
WORKDIR /app 
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -a -o sprawl . 

FROM alpine

RUN apk add leveldb

COPY --from=build /app/sprawl /sprawl

CMD ["/sprawl"]
