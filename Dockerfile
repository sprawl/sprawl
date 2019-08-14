FROM golang:latest as build

RUN mkdir /app 
COPY . /app/ 
WORKDIR /app 
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -a -o sprawl . 

FROM scratch

ENV SPRAWL_DATABASE_PATH /
ENV SPRAWL_API_PORT 1337

COPY --from=build /app/sprawl /sprawl

EXPOSE 1337

CMD ["/sprawl"]
