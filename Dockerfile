# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/bin/api ./cmd/api

# Final runtime image
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /root/
COPY --from=builder /app/bin/api /usr/local/bin/onduty-api
COPY db/migrations /root/db/migrations

ENV PORT=8080
ENV GIN_MODE=release

EXPOSE 8080

ENTRYPOINT ["onduty-api"]
