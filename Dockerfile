FROM golang:1.18-alpine
RUN apk add build-base

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
COPY . ./
RUN go mod download
