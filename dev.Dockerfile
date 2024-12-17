FROM golang:1.23-alpine3.21 as builder

RUN apk add --no-cache \
    build-base \
    libmediainfo-dev \
     openssl

WORKDIR /app

ENV GO111MODULE=on
ENV GOPATH /go


COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .


RUN go install github.com/air-verse/air@latest


CMD ["air", "-c", ".air.toml"]
