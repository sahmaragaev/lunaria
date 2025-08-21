# syntax=docker/dockerfile:1.4
FROM golang:1.23-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git ca-certificates tzdata
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -trimpath \
    -o lunaria-backend ./cmd

FROM scratch
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /app/lunaria-backend .
COPY --from=builder /app/config.yaml ./

USER 65534

EXPOSE 8080
EXPOSE 6060

ENV GIN_MODE=release

CMD ["./lunaria-backend", "server"]