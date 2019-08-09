FROM golang:latest as build

RUN mkdir /app 
COPY . /app/ 
WORKDIR /app 
RUN go build -o sprawl . 

FROM alpine

RUN apk add leveldb

COPY --from=build /app/sprawl /sprawl

CMD ["/sprawl"]
