FROM golang:1.18-buster

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
COPY . ./
RUN go mod download
RUN apt-get install build-essential
