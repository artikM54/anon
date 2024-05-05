FROM golang:latest

RUN mkdir /app
WORKDIR /app
ADD . /app

RUN go get github.com/cespare/reflex
RUN go install github.com/cespare/reflex

CMD reflex -r '\.go$' -s -- sh -c 'CONFIG_PATH="./" go run ./cmd/text_chat/main.go'
