# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o api-server ./cmd/api


# Runtime stage
FROM scratch

WORKDIR /

COPY --from=builder /app/api-server /api-server

EXPOSE 8080

ENTRYPOINT ["/api-server"]