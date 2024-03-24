ARG ALPINE_IMAGE=alpine:3.18
ARG SSO_BUILD_IMAGE=kitanoyoru/kitaDriveBot.sso:latest

FROM $SSO_BUILD_IMAGE AS build

FROM $ALPINE_IMAGE

USER root
RUN apk add --no-cache --upgrade libssl1.1 libcrypto1.1

RUN set -eux; \
    apk add --no-cache ca-certificates tzdata; \
    update-ca-certificates; \
    addgroup -g 1000 server; \
    adduser -HDu 1000 -G server server

COPY --from=build /usr/local/bin /usr/local/bin

USER 1000:1000

