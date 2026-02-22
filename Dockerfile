FROM golang:1.25-alpine AS builder

RUN apk update && apk add --no-cache git ca-certificates tzdata

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ARG VERSION=1.0.0
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -X main.Version=${VERSION}" \
    -o microservice \
    cmd/api/main.go

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o healthcheck \
    cmd/healthcheck/main.go

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

COPY --from=builder /app/microservice .
COPY --from=builder /app/healthcheck .
COPY --from=builder /app/config ./config

ENTRYPOINT ["./microservice"]