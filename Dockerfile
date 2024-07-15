#!/bin/sh
FROM golang:alpine
MAINTAINER Selly Margaretha Sudiyandi <sellymargaretha41@gmail.com>
RUN apk add --no-cache git
ENV GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64

RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o main .
CMD ["/app/main"]