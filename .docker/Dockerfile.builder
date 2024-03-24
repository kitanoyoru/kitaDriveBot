ARG BUILDER_IMAGE=golang:1.21-alpine
ARG ALPINE_IMAGE=alpine:3.18

FROM $BUILDER_IMAGE AS build

ARG GITHUB_USERNAME
ARG GITHUB_PAT
ARG GITHUB_HOST=github.com

ARG BUILD_PATH
ARG BUILD_RESULT

RUN apk update && apk add --no-cache git ca-certificates tzdata openssh-client && update-ca-certificates

RUN echo "machine ${GITHUB_HOST} login ${GITHUB_USERNAME} password ${GITHUB_PAT}" > ~/.netrc

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0
RUN go build -v -o /usr/local/bin/${BUILD_RESULT} ${BUILD_PATH} 


