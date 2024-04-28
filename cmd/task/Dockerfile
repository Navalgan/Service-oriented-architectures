FROM golang:latest

COPY ./ /app

WORKDIR /app/cmd/task

RUN go build main.go

ENTRYPOINT ./main
