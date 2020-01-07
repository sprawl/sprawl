FROM golang:latest as build

RUN mkdir /app
COPY . /app/
WORKDIR /app
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.configPath=" -a -o sprawl .

FROM scratch

COPY --from=build /app/sprawl /sprawl

ENV SPRAWL_DATABASE_PATH /home/sprawl/data
ENV SPRAWL_API_PORT 1337

CMD ["/sprawl"]
