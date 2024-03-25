ARG BUILDER_IMAGE=golang:1.21-alpine
ARG ALPINE_IMAGE=alpine:3.18

FROM $BUILDER_IMAGE AS build

ARG GITHUB_USERNAME
ARG GITHUB_PAT
ARG GITHUB_HOST=github.com

RUN apk update && apk add --no-cache git ca-certificates tzdata openssh-client && update-ca-certificates

RUN echo "machine ${GITHUB_HOST} login ${GITHUB_USERNAME} password ${GITHUB_PAT}" > ~/.netrc

WORKDIR /app

COPY . .

RUN go work sync 

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

RUN go build -gcflags='all=-N -l' -v -o /usr/local/bin/sso apps/sso/cmd/main.go 

