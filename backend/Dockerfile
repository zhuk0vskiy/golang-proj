FROM golang:1.22.4

RUN set -ex

USER root
WORKDIR /app

COPY src src
COPY go.mod go.mod

RUN go mod tidy
RUN cd src && go build main.go

CMD cd src && ./main