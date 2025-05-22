# Build stage
FROM golang:1.24.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build statically linked binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o turbogate ./cmd/main.go

# Final stage - using scratch (no OS deps)
FROM scratch

WORKDIR /app

COPY --from=builder /app/turbogate .

COPY --from=builder /app/config ./config

EXPOSE 10000

CMD ["./turbogate"]
