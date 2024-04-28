FROM golang:latest

COPY . /app

WORKDIR /app/task

RUN go build main.go

ENTRYPOINT ./main
