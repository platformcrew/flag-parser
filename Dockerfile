# Stage 1: Build the Go binary
FROM golang:1.24-alpine AS builder

WORKDIR /src

# Copy module files first for better layer caching
COPY go.mod go.sum* ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Build a fully static binary (no CGO, no libc dependency)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /flag-parser .

# Stage 2: Minimal runtime image
FROM alpine:3.21

RUN apk add --no-cache ca-certificates

COPY --from=builder /flag-parser /flag-parser

ENTRYPOINT ["/flag-parser"]
