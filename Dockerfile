# syntax=docker/dockerfile:1

FROM golang:latest

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
ADD handlers ./handlers

RUN go mod download

RUN go get smartmirror/handlers


COPY *.go ./


RUN go build -o /api-gateway

EXPOSE 8080

CMD [ "/api-gateway" ]