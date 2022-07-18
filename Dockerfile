FROM golang:1.18-alpine

WORKDIR /app
COPY go.mod ./
COPY . ./
RUN go mod download