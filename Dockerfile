# ----------------------------------
# Builder
# ----------------------------------
FROM golang:1.26.3-alpine AS builder

WORKDIR /app

# Dependencias necesarias
RUN apk add --no-cache git

# Descargar módulos
COPY go.mod go.sum ./
RUN go mod download

# Copiar código
COPY . .

# Compilar
RUN CGO_ENABLED=0 \
    go build \
    -trimpath \
    -ldflags="-s -w" \
    -o inariops \
    ./cmd/api

# ----------------------------------
# Runner
# ----------------------------------
FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/inariops .

EXPOSE 9142

CMD ["./inariops"]
