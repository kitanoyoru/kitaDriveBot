# syntax=docker/dockerfile:1

FROM golang:1.21-alpine AS builder

WORKDIR /src

RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/bot ./cmd/bot

FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /out/bot /usr/local/bin/bot

RUN adduser -D -H appuser
USER appuser

ENTRYPOINT ["bot"]
CMD ["run"]
