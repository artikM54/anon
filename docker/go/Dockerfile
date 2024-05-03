FROM golang:latest

RUN mkdir /app
WORKDIR /app
ADD . /app

CMD CONFIG_PATH="./" go run ./cmd/text_chat/main.go
