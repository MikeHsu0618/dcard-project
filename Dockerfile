# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.17-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build ./cmd/app/main.go

##
## Deploy
##
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY /resources /resources
COPY /.env /.env

COPY --from=build /app/main /main

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/main"]

