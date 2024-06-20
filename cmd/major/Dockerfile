FROM golang:latest

COPY . /app

WORKDIR /app/cmd/major

RUN go build main.go

ENTRYPOINT ./main --jwt-key my-sercret-jwt-key
