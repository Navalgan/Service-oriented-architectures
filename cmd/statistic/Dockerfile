FROM golang:latest

COPY . /app

WORKDIR /app/statistic

RUN go build main.go

ENTRYPOINT ./main
