FROM golang:latest

COPY . /app

WORKDIR /app/cmd/statistic

RUN go build main.go

ENTRYPOINT ./main
