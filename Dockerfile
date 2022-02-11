# syntax=docker/dockerfile:1

FROM golang:1.17-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go get github.com/codegangsta/gin
RUN go install github.com/codegangsta/gin

RUN go mod download

EXPOSE 8080