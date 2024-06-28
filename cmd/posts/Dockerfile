FROM golang:latest

COPY . /app

WORKDIR /app/cmd/posts

RUN go build main.go

ENTRYPOINT ./main
