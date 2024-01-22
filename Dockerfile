# syntax=docker/dockerfile:1

# build stage
FROM golang:1.21 AS build

WORKDIR /app
COPY . .

# test stage
FROM build AS test
WORKDIR /app
ENV TESTING=true
ENV AWS_ACCESS_KEY_ID=dummy
ENV AWS_SECRET_ACCESS_KEY=dummy
ENV AWS_DEFAULT_REGION=us-west-2
RUN go test -v ./dynamo/...